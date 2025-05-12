package server_test

import (
	"dqq/go/basic/socket/server"
	"testing"
)

func TestTcpServer(t *testing.T) {
	server.TcpServer()
}

func TestUdpServer(t *testing.T) {
	server.UdpServer()
}

func TestTcpLongConnection(t *testing.T) {
	server.TcpLongConnection()
}

func TestUdpLongConnection(t *testing.T) {
	server.UdpLongConnection()
}

func TestTcpStick(t *testing.T) {
	server.TcpStick()
}

func TestUdpConnectionCurrent(t *testing.T) {
	server.UdpConnectionCurrent()
}

func TestUdpRpcServer(t *testing.T) {
	server.UdpRpcServer()
}

// go test -v ./socket/server -run=^TestTcpServer$ -count=1
// go test -v ./socket/server -run=^TestUdpServer$ -count=1
// go test -v ./socket/server -run=^TestTcpLongConnection$ -count=1
// go test -v ./socket/server -run=^TestUdpLongConnection$ -count=1
// go test -v ./socket/server -run=^TestTcpStick$ -count=1
// go test -v ./socket/server -run=^TestUdpConnectionCurrent$ -count=1
// go test -v ./socket/server -run=^TestUdpRpcServer$ -count=1
