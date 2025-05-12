package test

import (
	"dqq/concurrency/raft"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"testing"
	"time"
)

func init() {
	fout, err := os.OpenFile("../raft.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	handler := slog.NewTextHandler(fout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func TestRaft(t *testing.T) {
	size := 7 //集群大小
	servers := make([]*raft.Server, 0, size)
	//启动各台raft server
	for i := 0; i < size; i++ {
		port := 7000 + i
		transporter := raft.NewHttpTransporter("/raft", 2*time.Second)
		server := raft.NewServer(fmt.Sprintf("http://127.0.0.1:%d", port), port, nil, transporter)
		server.Start(false)
		defer server.Stop()
		servers = append(servers, server)
	}
	//给每个server添加peer
	for _, server := range servers {
		for _, peer := range servers {
			server.AddPeer(&raft.Peer{Id: peer.Id, ConnectionString: peer.ConnectionString})
		}
	}
	//随机停掉一台server
	downServer := servers[rand.IntN(size)]
	downServer.Stop()

	//等待Leader产生
	var leader *raft.Server
	for { //如果没有产生leader，就一直执行这个for循环
		leader = getLeader(servers)
		if leader != nil { //如果集群内部自主选举出了leader
			fmt.Printf("leader elevationed %s\n", leader.ID()) //打印leader是谁
			break                                              //退出for循环
		}
		time.Sleep(1 * time.Second) //等待1秒，进入下一轮循环
	}

	//产生3倍的MaxLogEntriesPerHeartbeat
	for i := 0; i < 3*raft.MaxLogEntriesPerHeartbeat; i++ {
		leader.Do(raft.NoopCommand{})
		time.Sleep(50 * time.Millisecond) //休息50ms后，发送下一条log
	}

	//重启down server，它需要接收3倍的MaxLogEntriesPerHeartbeat，一次心跳发不完
	downServer.Start(true)
	//等所有日志复制完成
	time.Sleep(5 * time.Second)

	// 最终，集群中所有节点的LastLogIndex应该都是3*raft.MaxLogEntriesPerHeartbeat
}

// 并不是每个人心目中的leader都是一个，得找到票数最多的leader
func getLeader(servers []*raft.Server) *raft.Server {
	if len(servers) == 0 {
		return nil
	}
	countMap := make(map[string]int, len(servers))
	for _, server := range servers {
		countMap[server.LeaderID()] = countMap[server.LeaderID()] + 1
	}
	var leaderId string
	var maxVotes int = 1
	for k, v := range countMap {
		slog.Debug("vote count", k, v)
		if v > maxVotes {
			maxVotes = v
			leaderId = k
		}
	}
	for _, server := range servers {
		if leaderId == server.ID() {
			return server
		}
	}
	return nil
}

// go test -v ./raft/test -run=^TestRaft$ -count=1 -timeout=5m
