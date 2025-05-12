package main

import (
	"dqq/micro_service/my_rpc/transport"
	"fmt"
	"net"
	"time"
)

func main() {
	//跟tcp_client的唯一区别就是这行代码
	conn, err := net.DialTimeout("udp", "127.0.0.1:5678", 30*time.Minute) //一个conn绑定一个本地端口
	transport.CheckError(err)
	defer conn.Close()

	for i := 0; i < 3; i++ {
		requestBytes := []byte("china")
		_, err = conn.Write(requestBytes)
		transport.CheckError(err)
		fmt.Printf("write request %s\n", string(requestBytes))
	}
}

// go run ./micro_service/my_rpc/transport/client
