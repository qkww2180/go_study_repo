package server

import (
	"bytes"
	transport "dqq/go/basic/socket"
	"io"
	"log"
	"net"
	"time"
)

func TcpServer() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	log.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	log.Printf("establish connection to client %s\n", conn.RemoteAddr().String())
	conn.SetReadDeadline(time.Now().Add(10 * time.Second)) //设置一个读的期限，超过这个期限再调Read()就会发生error。默认是60s内可Read()。
	defer conn.Close()                                     //关闭连接

	request := make([]byte, 256)
	n, err := conn.Read(request) //跟读文件类似。可能会因为超时而导致error(之前设置了ReadDeadline)
	transport.CheckError(err)
	log.Printf("receive %s\n", string(request[:n]))
} // TCP需要先启动server，再启动client，否则client连接不上server

func TcpLongConnection() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	log.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	log.Printf("establish connection to client %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	// time.Sleep(5 * time.Second) //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256)
	for { //长连接
		conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) //每次都要续命
		n, err := conn.Read(request)                          //如果流里没数据，Read()会阻塞。对方close后，这里会遇到EOF
		transport.CheckError(err)
		log.Printf("receive %s\n", string(request[:n]))
	}
}

// TCP粘包问题。解决办法：
//
// 1. 通信双方约定好，消息为固定长度。
//
// 2. 双方约定好消息之间的分隔符。
//
// 3. 定义好消息的序列化方式，约定好序列化之后的第几个字节是消息长度。
func TcpStick() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	transport.CheckError(err)
	log.Println("waiting for client connection ......")
	conn, err := listener.Accept()
	transport.CheckError(err)
	log.Printf("establish connection to client %s\n", conn.RemoteAddr().String())
	defer conn.Close()

	time.Sleep(5 * time.Second) //故意多sleep一会儿，让client多发几条消息过来
	request := make([]byte, 256)
	buffer := bytes.Buffer{}
	for { //只要client不关闭连接，server就得随时待命
		n, err := conn.Read(request) //TCP是面向字节流的，一次Read到的数据可能包含了多个报文，也可能只包含了半个报文，一条报文在什么地方结束需要通信双方事先约定好

		// bufio的ReadBytes只支持1个byte作为分隔符，而我们的分隔符是多个byte，所以只能自行实现stream的分隔逻辑

		if err != nil {
			if err == io.EOF { //对方关闭了连接
				if buffer.Len() > 0 {
					log.Println(buffer.String())
				}
			} else {
				log.Printf("read error %v", err)
			}
			break
		}
		log.Printf("receive request %s\n", string(request[:n]))
		data := request[:n] //如果分隔符刚好横跨2个data，则buffer里会包含2个报文，中间含一个分隔符
		for {               // data里可能包含多个分割符，所以需要循环解析
			pos := bytes.Index(data, transport.MAGIC) //约定好用MAGIC当分割符
			if pos >= 0 {
				if pos == 0 { //data的起始位置刚好就是分割符
					if buffer.Len() > 0 {
						useBuffer(&buffer)
					}
				} else if pos > 0 {
					buffer.Write(data[:pos]) //分割符之前的内容追加到buffer里
					useBuffer(&buffer)
				}
				data = data[pos+len(transport.MAGIC):] //游标往前走
			} else {
				buffer.Write(data) //data里没包含分割符，则把data全部追加到buffer里
				break
			}
		}
	}
}

func useBuffer(buffer *bytes.Buffer) {
	defer buffer.Reset() //清空buffer
	data := buffer.Bytes()
	pos := bytes.Index(data, transport.MAGIC)
	if pos < 0 {
		log.Println(buffer.String()) //把buffer里的内容输出
	} else {
		//buffer里可能包含一个或两个报文
		log.Println(string(data[:pos]))
		log.Println(string(data[pos+len(transport.MAGIC):]))
	}
}
