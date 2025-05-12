package grpcconnpool

import (
	"log"
	"time"

	"google.golang.org/grpc"
)

type connection struct {
	conn     *grpc.ClientConn
	deadline time.Time
}

// 重新建立连接
func (connection *connection) refresh(cfg ClientConfig, opts ...grpc.DialOption) {
	if connection == nil {
		return
	}
	if connection.conn != nil {
		connection.conn.Close()
	}
	conn, err := grpc.Dial(cfg.serverAddr, opts...)
	if err != nil {
		log.Printf("connect to %s failed: %v", cfg.serverAddr, err)
		return
	}
	connection.conn = conn
	connection.deadline = time.Now().Add(cfg.GenLifeTime())
	// log.Printf("connect to %s", cfg.serverAddr)
}

// 判断是否需要重新建立连接
func (connection *connection) shouldRefresh() bool {
	if connection.conn == nil {
		return true
	}
	if !isConnectionHealthy(connection.conn) {
		return true
	}
	if time.Now().After(connection.deadline) {
		return true
	}
	return false
}
