package main

import (
	"dqq/micro_service/my_rpc"
	"dqq/micro_service/my_rpc/serialization"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/rs/xid"
)

type EchoClient struct {
	conn          net.Conn
	serializer    serialization.Serializer
	requestBuffer sync.Map
}

func NewEchoClient(serverIp string, serverPort int, serializer serialization.Serializer) *EchoClient {
	conn, err := net.DialTimeout("udp", serverIp+":"+strconv.Itoa(serverPort), 30*time.Minute) // UDP是面向报文的，TCP是面向字节流的，用UDP实现方便一些，若用TCP实现需要自行按魔数对字节流进行分隔
	if err != nil {
		log.Println(err)
		return nil
	}
	client := EchoClient{
		conn:          conn,
		requestBuffer: sync.Map{},
		serializer:    serializer,
	}
	go client.receive()
	return &client
}

func (client *EchoClient) receive() {
	for {
		bs := make([]byte, 4096)
		n, err := client.conn.Read(bs)
		if err != nil {
			log.Println(err)
		} else {
			var resp my_rpc.RpcData
			if err = client.serializer.Unmarshal(bs[:n], &resp); err != nil {
				log.Println(err)
			} else {
				if v, exists := client.requestBuffer.Load(resp.Id); exists {
					client.requestBuffer.Delete(resp.Id)
					ch := v.(chan my_rpc.RpcData)
					ch <- resp
				} else {
					log.Printf("request buffer里没有id %s", resp.Id)
				}
			}
		}
	}
}

// 因为在Call函数内要修改EchoClient的成员requestBuffer，所以必须传EchoClient的指针
func (client *EchoClient) Call(request *my_rpc.RpcData) *my_rpc.RpcData {
	id := xid.New().String() //分布式唯一id生成器
	request.Id = id
	data, err := client.serializer.Marshal(*request) //序列化request
	if err != nil {
		log.Println(err)
		return nil
	}

	streamCh := make(chan my_rpc.RpcData, 1)
	client.requestBuffer.Store(id, streamCh)
	_, err = client.conn.Write(data) //发送请求
	if err != nil {
		log.Println(err)
		return nil
	}
	resp := <-streamCh //阻塞，直到能从channel里取出响应
	return &resp
}

func main() {
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)

	for p := 0; p < P; p++ {
		go func() {
			defer wg.Done()
			client := NewEchoClient("127.0.0.1", 5678, my_rpc.Serializer)
			for i := 0; i < 10; i++ {
				a := rand.Int()
				request := my_rpc.RpcData{A: a}
				resp := client.Call(&request) //执行RPC调用
				if resp == nil {
					log.Printf("rpc调用失败")
					continue
				}
				if resp.A != a {
					log.Printf("rpc调用结果不符合预期，%d!=%d", resp.A, a)
				}
			}
		}()
	}

	wg.Wait()
}

// 先启server，再启client
// go run ./micro_service/my_rpc/client
