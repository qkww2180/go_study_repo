package main

import (
	"bufio"
	"dqq/util"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type LogCollector struct {
	conn   *net.UDPConn
	fout   *os.File
	writer *bufio.Writer
}

func NewLogCollector(port int, file string) (*LogCollector, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	//此conn不支持并发使用
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	fout, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //注意：日志文件一般都是append模式
	if err != nil {
		return nil, err
	}
	writer := bufio.NewWriter(fout)
	collector := &LogCollector{
		conn:   conn,
		fout:   fout,
		writer: writer,
	}
	return collector, nil
}

// 能过UDP接收Log，写入buffer
func (c LogCollector) Receive() {
	content := make([]byte, 4<<20) //一条日志不会超过4M，如果真超过了，ReadFromUDP时会报错
	for {                          //Receive 会一直阻塞
		read_len, remoteAddr, err := c.conn.ReadFromUDP(content)
		if err != nil { //如果content的长度小于报文的大小，会发生error
			fmt.Printf("receive log failed: %v\n", err)
		} else {
			log := "[" + remoteAddr.IP.String() + "] " + string(content[:read_len]) //标记这条日志来自哪台服务器
			// c.buffer <- log
			if _, err := c.writer.Write([]byte(log)); err != nil {
				fmt.Printf("write log <%s> to file fail: %v\n", log, err)
			} else {
				c.writer.WriteString("\n")
			}
		}
	}
}

func (c LogCollector) Close() {
	c.conn.Close() //关闭socket连接
	c.writer.Flush()
	c.fout.Close()
}

// 需要监听信号2(Ctrl+C - SIGINT)和15(SIGTERM)，当收到信号时关闭LogCollector
func listenSignal(collector *LogCollector) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
	sig := <-c                                        //阻塞，直到信号的到来
	fmt.Printf("receive signal %s\n", sig.String())
	if collector != nil {
		collector.Close()
	}
	os.Exit(0) //进程退出
}

func main() {
	collector, err := NewLogCollector(5678, util.RootPath+"log/collector.log")
	if err != nil {
		panic(err)
	}
	go listenSignal(collector)
	collector.Receive()
}

// go run .\j_io\logger\udp_logger\collector\
