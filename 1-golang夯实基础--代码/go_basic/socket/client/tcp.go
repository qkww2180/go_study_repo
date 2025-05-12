package client

import (
	transport "dqq/go/basic/socket"
	"log"
	"net"
)

func connect2TcpServer(serverAddr string) *net.TCPConn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	transport.CheckError(err)
	log.Printf("ip %s port %d\n", tcpAddr.IP.String(), tcpAddr.Port)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	transport.CheckError(err)
	log.Printf("establish connection to server %s myself %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String()) //操作系统会随机给客户端分配一个49152~65535上的端口号
	return conn
}

func sendTcpServer(conn net.Conn) {
	// n, err := conn.Write([]byte("hello")) //跟写文件类似
	n, err := conn.Write(append([]byte("hello"), transport.MAGIC...))
	transport.CheckError(err)
	log.Printf("send %d bytes\n", n)
}

func TcpClient() {
	conn := connect2TcpServer("127.0.0.1:5678")

	sendTcpServer(conn)
	conn.Close()
	log.Println("close connection")
}

func TcpLongConnection() {
	conn := connect2TcpServer("127.0.0.1:5678")

	for i := 0; i < 3; i++ {
		sendTcpServer(conn)
	}
	conn.Close()
	log.Println("close connection")
}

func TcpStick() {
	conn := connect2TcpServer("127.0.0.1:5678")

	for i := 0; i < 3; i++ {
		sendTcpServer(conn)
	}
	conn.Close()
	log.Println("close connection")
}
