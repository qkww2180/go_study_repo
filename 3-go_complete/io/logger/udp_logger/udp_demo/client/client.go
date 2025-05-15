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
	//跟tcp_client的唯一区别就是这行代码
	conn, err := net.DialTimeout("udp", ip+":"+strconv.Itoa(port), 30*time.Minute) //一个conn绑定一个本地端口
	if err != nil {
		fmt.Printf("dial to server failed: %s", err)
		return
	}
	defer conn.Close()

	requestBytes := []byte("china")
	_, err = conn.Write(requestBytes)
	if err != nil {
		fmt.Printf("write to server failed: %s", err)
		return
	}
	fmt.Printf("write request %s\n", string(requestBytes))
	responseBytes := make([]byte, 256) //初始化后byte数组每个元素都是0
	read_len, err := conn.Read(responseBytes)
	if err != nil {
		fmt.Printf("read from server failed: %s", err)
		return
	}
	fmt.Printf("receive response: %s\n", string(responseBytes[:read_len]))
}

// go run .\j_io\logger\udp_logger\udp_demo\client\client.go
