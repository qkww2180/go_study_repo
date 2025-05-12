package main

import (
	"bytes"
	"dqq/micro_service/my_rpc/transport"
	"fmt"
	"io"
	"net"
	"time"
)

// TCP是面向字节流的
func main1() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	fmt.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	fmt.Printf("establish connection to client %s\n", conn.RemoteAddr().String()) //操作系统会随机给客户端分配一个49152~65535上的端口号

	time.Sleep(5 * time.Second)  //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256) //设定一个最大长度，防止flood attack
	buffer := bytes.Buffer{}
	for { //只要client不关闭连接，server就得随时待命
		n, err := conn.Read(request) //TCP是面向字节流的，一次Read到的数据可能包含了多个报文，也可能只包含了半个报文，一条报文在什么地方结束需要通信双方事先约定好
		if err == io.EOF {           //对方关闭了连接
			fmt.Println(buffer.String())
			break
		}
		fmt.Printf("receive request %s\n", string(request[:n]))
		data := request[:n]
		for {
			pos := bytes.Index(data, transport.MAGIC) //约定好用MAGIC当分割符
			if pos >= 0 {
				if pos == 0 {
					if buffer.Len() > 0 {
						fmt.Println(buffer.String())
						buffer.Reset()
					}
				} else if pos > 0 {
					buffer.Write(data[:pos])
					fmt.Println(buffer.String())
					buffer.Reset()
				}
				data = data[pos+len(transport.MAGIC):]
			} else {
				buffer.Write(data)
				break
			}
		}
	}
	conn.Close()
}

// go run ./micro_service/my_rpc/transport/server
