package client

import (
	transport "dqq/go/basic/socket"
	"encoding/json"
	"log"
	rand "math/rand/v2"
	"net"
	"sync"
	"time"
)

func connect2UdpServer(serverAddr string) net.Conn {
	//跟tcp_client的唯一区别就是这行代码
	conn, err := net.DialTimeout("udp", serverAddr, 30*time.Minute) //一个conn绑定一个本地端口
	transport.CheckError(err)
	log.Printf("establish connection to server %s myself %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String()) //操作系统会随机给客户端分配一个49152~65535上的端口号
	return conn
}

func sendUdpServer(conn net.Conn) {
	n, err := conn.Write([]byte("hello")) //即使Server还启动，建立连接和发送数据都不会返回error，Server启动后也收不到这个数据
	transport.CheckError(err)
	log.Printf("send %d bytes\n", n)
}

func UdpClient() {
	conn := connect2UdpServer("127.0.0.1:5678")

	sendUdpServer(conn)
	conn.Close()
	log.Println("close connection")
}

func UdpLongConnection() {
	conn := connect2UdpServer("127.0.0.1:5678")

	for i := 0; i < 3; i++ {
		sendUdpServer(conn)
	}

	time.Sleep(70 * time.Second)
	sendUdpServer(conn)
	conn.Close()
	log.Println("close connection")
}

// Client端，并发使用udp连接
func UdpConnectionCurrent() {
	conn := connect2UdpServer("127.0.0.1:5678")

	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			sendUdpServer(conn)
		}()
	}
	wg.Wait()
}

func UdpRpcClient() {
	const P = 500 // 模拟500个client
	const C = 10  //每个client发起10次请求，然后关闭连接
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			conn := connect2UdpServer("127.0.0.1:5678")
			for j := 0; j < C; j++ {
				request := transport.AddRequest{
					RequestId: rand.Int(),
					A:         int(rand.Int32()) % 100,
					B:         int(rand.Int32()) % 100,
				}
				bs, err := json.Marshal(request)
				if err != nil {
					log.Printf("marshal request failed: %s", err)
					continue
				}
				if _, err := conn.Write(bs); err == nil {
					log.Printf("send request, id %d a %d b %d", request.RequestId, request.A, request.B)
				}
			}

			buffer := make([]byte, 256)
			for j := 0; j < C; j++ {
				if n, err := conn.Read(buffer); err == nil {
					var response transport.AddResponse
					err = json.Unmarshal(buffer[:n], &response)
					if err == nil {
						log.Printf("receive response, id %d sum %d", response.RequestId, response.Sum)
					} else {
						log.Printf("unmarshal response failed: %s", err)
					}
				} else {
					log.Printf("read response failed: %s", err)
				}
			}
			conn.Close() //每个client发起10次请求，然后关闭连接
		}()
	}
	wg.Wait()
}
