package grpcconnpool

import (
	"math"
	rand "math/rand/v2"
	"time"
)

type ClientConfig struct {
	serverAddr                  string        //服务端地址
	connectionPoolSize          int           //连接池大小
	connectionLifeTime          time.Duration //连接的寿命，到期后强制废弃，建立新连接。<=0表示没有到期时间
	connectionLifeTimeDeviation time.Duration //为避免所有连接同时到期，给寿命设置一个随机波动
}

/**
Set系列函数。Builder模式
*/

func (cc *ClientConfig) WithServerAddr(arg string) *ClientConfig {
	cc.serverAddr = arg
	return cc
}

func (cc *ClientConfig) WithConnectionPoolSize(arg int) *ClientConfig {
	cc.connectionPoolSize = arg
	return cc
}

func (cc *ClientConfig) WithConnectionLifeTime(arg time.Duration) *ClientConfig {
	cc.connectionLifeTime = arg
	return cc
}

func (cc *ClientConfig) WithConnectionLifeTimeDeviation(arg time.Duration) *ClientConfig {
	cc.connectionLifeTimeDeviation = arg
	return cc
}

// 根据随机波动，生成一个寿命
func (cc *ClientConfig) GenLifeTime() time.Duration {
	if cc.connectionLifeTime > 0 {
		return cc.connectionLifeTime + time.Duration(int64(rand.Float64()*float64(cc.connectionLifeTimeDeviation)))
	} else {
		return time.Duration(math.MaxInt64 - time.Now().Unix())
	}
}
