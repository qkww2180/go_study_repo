package raft

import (
	"time"
)

const (
	HeartBeatInterval         = 20 * time.Millisecond  //心跳间隔。由于收到log entry后还得写磁盘，所以这个时间也不能太短，建议0.5-20ms
	ElectionTimeout           = 150 * time.Millisecond //follower多久收不到leader的心跳，就变为candidate。这个时间要明显大于heartBeatInterval
	LeaderChangeTimeout       = 500 * time.Millisecond //选主过程的超时，建议10-500ms。这个时间要明显大于heartBeatInterval
	MaxLogEntriesPerHeartbeat = 100                    //每次心跳最多发送多少条日志，太多了会超出rpc data size的限制，也会让leader阻塞更长时间
)

// 检查全局参数的合法性
func init() {
	if MaxLogEntriesPerHeartbeat < 1 {
		panic("maxLogEntriesPerHeartbeat must more than one")
	}
	if 2*HeartBeatInterval >= ElectionTimeout { //ElectionTimeout要显著大于HeartBeatInterval
		panic("heartBeatInterval should much less than electionTimeout")
	}
	if 2*HeartBeatInterval >= LeaderChangeTimeout { //LeaderChangeTimeout要显著大于HeartBeatInterval
		panic("heartBeatInterval should much less than electionTimeout")
	}
}
