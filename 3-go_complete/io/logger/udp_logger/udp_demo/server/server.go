package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// 收发简单的字符串消息
func main() {
	ip := "127.0.0.1" //ip换成0.0.0.0和空字符串试试
	port := 5656
	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("resolve address failed: %s", err)
		return
	}
	//conn 不支持并发使用
	conn, err := net.ListenUDP("udp", udpAddr) //UDP不需要创建连接，所以不需要像TCP那样通过Accept()创建连接，这里的conn是个假连接
	if err != nil {
		fmt.Printf("listen failed: %s", err)
		return
	}
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) //设置读超时
	defer conn.Close()

	requestBytes := make([]byte, 256)                           //初始化后byte数组每个元素都是0
	read_len, remoteAddr, err := conn.ReadFromUDP(requestBytes) //一个conn可以对应多个client，ReadFrom可以返回是哪个。UDP是面向报文的，一次read刚好读好一条数据。如果由于buffer小导致读不完一条完整的消息（也取不到remoteAddr），则后半部分下次read也读不到了
	if err != nil {
		fmt.Printf("read from client error: %s\n", err.Error())
		return //到达deadline后，退出for循环，关闭连接。client再用这个连接读写会发生错误
	}
	fmt.Printf("receive request %s from %s\n", string(requestBytes[:read_len]), remoteAddr.String()) //[]byte转string时，0后面的会自动被截掉

	response := "hello " + string(requestBytes[:read_len])
	_, err = conn.WriteToUDP([]byte(response), remoteAddr) //由于UDP conn支持多对多通信，所以通信对方可能有多个EndPoint，通过WriteTo指定要写给哪个EndPoint
	if err != nil {
		fmt.Printf("write to client failed: %s", err)
		return
	}
	fmt.Printf("write response %s to %s\n", response, remoteAddr.String())
}

// go run .\io\logger\udp_logger\udp_demo\server\server.go
