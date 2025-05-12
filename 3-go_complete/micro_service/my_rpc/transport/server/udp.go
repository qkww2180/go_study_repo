package main

import (
	"dqq/micro_service/my_rpc/transport"
	"fmt"
	"net"
	"time"
)

// UDP是面向报文的，一次read刚好读好一条数据。如果由于buffer小导致读不完一条完整的消息，则后半部分下次read也读不到了
func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5678")
	transport.CheckError(err)
	conn, err := net.ListenUDP("udp", udpAddr) //UDP不需要创建连接，所以不需要像TCP那样通过Accept()创建连接，这里的conn是个假连接
	transport.CheckError(err)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	defer conn.Close()

	time.Sleep(5 * time.Second) //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256)
	for {
		read_len, remoteAddr, err := conn.ReadFromUDP(request) //UDP是面向报文的，一次read刚好读好一条数据。如果由于buffer小导致读不完一条完整的消息（也取不到remoteAddr），则后半部分下次read也读不到了
		transport.CheckError(err)
		fmt.Printf("receive request %s from %s\n", string(request[:read_len]), remoteAddr.String())
	}
}

// go run ./micro_service/my_rpc/transport/server
