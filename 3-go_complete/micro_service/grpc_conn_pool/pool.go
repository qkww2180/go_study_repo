package grpcconnpool

import (
	"log"

	"google.golang.org/grpc"
)

type GrpcClientPool struct {
	cfg          ClientConfig
	conns        []*connection
	dialOpts     []grpc.DialOption
	loadBalancer LoadBalancer //策略模式。把策略的实现转移到其他类里面去
}

// 获取一个可用的连接。如果没有可用的会返回nil
func (pool *GrpcClientPool) GetConn() *grpc.ClientConn {
	idx := pool.loadBalancer.Select(len(pool.conns))
	conn := pool.conns[idx] //先随机选一个连接

	if conn.shouldRefresh() {
		conn.refresh(pool.cfg, pool.dialOpts...)
	}

	//如果选中的连接不可用，则依次往后找一个能用的连接
	if conn != nil && isConnectionHealthy(conn.conn) {
		// log.Printf("select %d connection", idx)
		return conn.conn
	} else {
		return pool.getNextHealthyClient(idx, len(pool.conns))
	}
}

func (pool *GrpcClientPool) getNextHealthyClient(currIndex, max int) *grpc.ClientConn {
	i := currIndex
	for {
		i = (i + 1) % max
		if i == currIndex { //轮了一遍，没找到可用的连接
			return nil
		}
		if pool.conns[i] != nil && isConnectionHealthy(pool.conns[i].conn) {
			// log.Printf("select %d connection", i)
			return pool.conns[i].conn
		}
	}
}

func NewGrpcClientPool(cfg ClientConfig, dialOpts []grpc.DialOption, loadBalancer LoadBalancer) *GrpcClientPool {
	if cfg.connectionPoolSize <= 0 {
		log.Printf("invalid pool size %d", cfg.connectionPoolSize)
		return nil
	}
	if len(cfg.serverAddr) == 0 {
		log.Printf("please indicate server address")
		return nil
	}
	conns := make([]*connection, 0, cfg.connectionPoolSize)
	for i := 0; i < cfg.connectionPoolSize; i++ {
		conn := new(connection)
		conn.refresh(cfg, dialOpts...)
		conns = append(conns, conn)
	}

	return &GrpcClientPool{
		cfg:          cfg,
		conns:        conns,
		dialOpts:     dialOpts,
		loadBalancer: loadBalancer,
	}
}
