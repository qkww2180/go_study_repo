package main

import (
	transport "dqq/go/basic/e_socket"
	"dqq/go/basic/f_encryption"
	"fmt"
	"log"
	"net"
	"time"
)

func init() {
	f_encryption.ReadRSAKey("./z_data/rsa_public_key.pem", "./z_data/rsa_private_key.pem")
}

func main() {
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

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer) //读取RSA加密之后的客户端的AES key
	transport.CheckError(err)
	aesKey, err := f_encryption.RsaDecrypt(buffer[:n]) //RSA解密
	transport.CheckError(err)
	// fmt.Printf("aes key %s\n", string(aesKey))
	conn.Write([]byte("I receive aes key"))

	key := [16]byte{}
	if len(aesKey) != 16 {
		panic(len(aesKey))
	}
	for i := 0; i < 16; i++ {
		key[i] = aesKey[i]
	}

	n, err = conn.Read(buffer) //读取AES加密之后的机密数据
	transport.CheckError(err)
	plain, err := f_encryption.AesDecrypt(string(buffer[:n]), key)
	transport.CheckError(err)
	fmt.Println(plain)
}

// go run ./f_encryption/tls/server
