package main

import (
	"dqq/micro_service/my_rpc/transport"
	"fmt"
	"net"
)

func main1() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	fmt.Printf("ip %s port %d\n", tcpAddr.IP.String(), tcpAddr.Port)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	transport.CheckError(err)
	fmt.Printf("establish connection to server %s\n", conn.RemoteAddr().String())

	//连发3条报文
	for i := 0; i < 3; i++ {
		n, err := conn.Write(append([]byte("china"), transport.MAGIC...)) //跟server约定好用MAGIC做为分隔符
		transport.CheckError(err)
		fmt.Printf("write request %d bytes\n", n)
	}

	conn.Close()
	fmt.Println("close connection")
}

// go run ./micro_service/my_rpc/transport/client
