package main

import (
	"dqq/go/basic/encryption"
	transport "dqq/go/basic/socket"
	"log"
	"net"
)

func init() {
	encryption.ReadRSAKey("./data/rsa_public_key.pem", "./data/rsa_private_key.pem")
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:5678")
	transport.CheckError(err)
	log.Printf("ip %s port %d\n", tcpAddr.IP.String(), tcpAddr.Port)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	transport.CheckError(err)
	log.Printf("establish connection to server %s myself %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
	defer conn.Close()

	aesKey := []byte("ir489u58ir489u54")
	decrypted, err := encryption.RsaEncrypt(aesKey)
	transport.CheckError(err)
	_, err = conn.Write(decrypted)
	transport.CheckError(err)
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer) //先阻塞一下，确保第一个阶段执行完毕，再发下一条数据，避免TCP的粘包问题
	transport.CheckError(err)

	key := [16]byte{}
	if len(aesKey) != 16 {
		panic(len(aesKey))
	}
	for i := 0; i < 16; i++ {
		key[i] = aesKey[i]
	}

	plain := "明月多情应笑我"
	s, err := encryption.AesEncrypt(plain, key)
	transport.CheckError(err)
	_, err = conn.Write([]byte(s))
	transport.CheckError(err)
}

// go run ./encryption/tls/client
