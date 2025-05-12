package main

import (
	"dqq/micro_service/my_rpc"
	"dqq/micro_service/my_rpc/serialization"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type EchoServer struct {
	conn       *net.UDPConn //UDP是面向报文的，TCP是面向字节流的，用UDP实现方便一些，若用TCP实现需要自行按魔数对字节流进行分隔
	serializer serialization.Serializer
}

func NewEchoServer(port int, serializer serialization.Serializer) *EchoServer {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return nil
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println("waiting for client connection ......")
	return &EchoServer{
		conn:       conn,
		serializer: serializer,
	}

}

func (server EchoServer) Serve() {
	defer server.conn.Close()
	for {
		request := make([]byte, 4096)                          //设定一个最大长度，防止flood attack
		n, remoteAddr, err := server.conn.ReadFromUDP(request) //读取请求
		if err != nil {
			log.Println(err)
			return
		}
		go server.handle(request[:n], remoteAddr) //并发处理每一个请求，提高吞吐量
	}
}

func (server EchoServer) handle(request []byte, remoteAddr *net.UDPAddr) {
	var data my_rpc.RpcData
	err := server.serializer.Unmarshal(request, &data) //反序列化request
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("receive request %+v from %s\n", data, remoteAddr.String())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(500))) //随机休息一段时间，先请求的不一定先拿到响应
	stream, err := server.serializer.Marshal(data)               //直接把接收到的内容再返回给对方
	if err != nil {
		log.Println(err)
		return
	}
	_, err = server.conn.WriteToUDP(stream, remoteAddr) //发送响应
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("send response  to %s\n", remoteAddr.String())
}

func main() {
	server := NewEchoServer(5678, my_rpc.Serializer)
	server.Serve()
}

// 先启server，再启client
// go run ./micro_service/my_rpc/server
