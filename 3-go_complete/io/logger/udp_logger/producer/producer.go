package main

import (
	"fmt"
	"net"
	"sync"
)

type LogProducer struct {
	buffer  chan string
	conn    net.Conn
	flushed chan struct{}
}

func NewLogProducer(collectorAddr string, bufferSize int) (*LogProducer, error) {
	if conn, err := net.Dial("udp", collectorAddr); err == nil { //创建UDP连接
		producer := &LogProducer{
			buffer:  make(chan string, bufferSize),
			conn:    conn,
			flushed: make(chan struct{}, 1),
		}
		go producer.daemonSend()
		return producer, nil
	} else {
		return nil, err
	}
}

func (p LogProducer) daemonSend() {
	for {
		if log, ok := <-p.buffer; ok {
			if _, err := p.conn.Write([]byte(log)); err != nil {
				fmt.Printf("send log <%s> fail: %v\n", log, err)
			}
		} else {
			p.flushed <- struct{}{} //buffer管道已关闭，且已被清空
			break
		}
	}
}

func (p LogProducer) Send(log string) (err error) {
	defer func() {
		if obj := recover(); obj != nil {
			err = fmt.Errorf("%v", obj)
		}
	}()
	p.buffer <- log //异步发送
	return
}

func (p LogProducer) Close() {
	close(p.buffer) //关闭buffer管道，不允许再往里面发送内容
	<-p.flushed     //等待buffer管道被清空
	p.conn.Close()  //关闭socket连接
}

func main() {
	producer, err := NewLogProducer("127.0.0.1:5678", 100)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				if err := producer.Send("好雨知时节，当春乃发生。随风潜入夜，润物细无声。野径云俱黑，江船火独明。晓看红湿处，花重锦官城。"); err != nil {
					fmt.Printf("发送日志失败:%s\n", err.Error())
				}
			}
		}()
	}
	wg.Wait()
}

// go run .\io\logger\udp_logger\producer\
