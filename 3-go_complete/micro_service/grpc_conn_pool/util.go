package grpcconnpool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func isConnectionHealthy(conn *grpc.ClientConn) bool {
	if conn == nil {
		return false
	}
	return conn.GetState() == connectivity.Connecting || conn.GetState() == connectivity.Ready
}
