package raft

import (
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/rs/xid"
)

type Peer struct {
	Id               string
	ConnectionString string
}

// raft节点
type Server struct {
	sync.RWMutex
	Peer
	port         int
	term         int64
	leaderId     string
	votedFor     string
	state        State
	log          *Log
	peers        []*Peer          //集群中的其他节点，要跟它们进行RPC通信
	prevLogIndex map[string]int64 //维护每个peer同步到了哪条日志
	fsm          FSM
	transporter  Transporter

	shutdownCh   chan struct{}
	rpcCh        chan RPC
	routineGroup sync.WaitGroup
}

func NewServer(connString string, port int, fsm FSM, transporter Transporter) *Server {
	server := &Server{
		port:         port,
		peers:        make([]*Peer, 0, 8),
		prevLogIndex: make(map[string]int64, 8),
		fsm:          fsm,
		transporter:  transporter,
		shutdownCh:   make(chan struct{}, 1),
		rpcCh:        make(chan RPC, 100),
	}
	server.log = NewLog(server)
	server.Id = xid.New().String()
	server.ConnectionString = connString
	return server
}

// LeaderID 当前leader的id
func (s *Server) LeaderID() string {
	return s.leaderId
}

func (s *Server) ID() string {
	return s.Id
}

// 超过一半即为大多数
func (s *Server) QuorumSize() int {
	return (len(s.peers)+1)/2 + 1
}

// 添加一个通信节点
func (s *Server) AddPeer(peer *Peer) {
	if peer == nil {
		return
	}
	if len(peer.Id) == 0 || len(peer.ConnectionString) == 0 {
		return
	}
	if peer.Id == s.Id { //把自己排除掉
		return
	}
	s.peers = append(s.peers, peer)
}

// 升级term，把votedFor清空，重置leaderId，把自己降为follower
func (s *Server) upgradeTerm(term int64, leaderId string) {
	s.term = term
	s.votedFor = "" //把votedFor清空
	s.leaderId = leaderId
	s.SetState(Follower)
}

func (s *Server) GetState() State {
	s.RLock()
	defer s.RUnlock()
	return s.state
}

func (s *Server) SetState(state State) {
	s.Lock()
	defer s.Unlock()
	s.state = state
	if state == Leader {
		s.leaderId = s.Id
	}
}

func (s *Server) Start(restart bool) {
	if !restart { //仅第一次启动时需要启动http server，Stop()时不会停掉http server
		slog.Info("start raft server", "id", s.Id)
		go s.transporter.Start(s.port, s) //启动http server
	} else {
		slog.Info("restart raft server", "id", s.Id)
		s.state = Follower
		s.shutdownCh = make(chan struct{}, 1)
		s.rpcCh = make(chan RPC, 100)
		s.routineGroup = sync.WaitGroup{}
	}
	go s.print()
	s.routineGroup.Add(1)
	go func() {
		defer s.routineGroup.Done()
		for s.GetState() != Stopped {
			switch s.GetState() {
			case Follower:
				s.FollowerLoop()
			case Candidate:
				s.CandidateLoop()
			case Leader:
				s.LeaderLoop()
			}
		}
	}()
}

// 打印集群的log index和commit index
func (s *Server) print() {
	ticker := time.NewTicker(3 * time.Second)
	for {
		<-ticker.C
		if s.GetState() == Leader {
			fmt.Println("raft cluster info", "leader", s.Id, "log index", s.log.LastIndex(), "commit index", s.log.CommitIndex())
			for _, peer := range s.peers {
				prevLogIndex := s.prevLogIndex[peer.Id]
				fmt.Println("raft cluster info", "follower", peer.Id, "log index", prevLogIndex)
			}
		}
	}
}

func (s *Server) Stop() {
	if s.GetState() == Stopped {
		return
	}
	s.state = Stopped
	s.shutdownCh <- struct{}{}
	close(s.shutdownCh)
	close(s.rpcCh)
	s.routineGroup.Wait() //等所有的异步任务结束
	slog.Info("server shutdown", "id", s.Id)
}

func (s *Server) FollowerLoop() {
	slog.Info("server run as follower", "id", s.Id, "leader id", s.leaderId)
	electionTimer := randomTimeout(ElectionTimeout)

	for s.state == Follower {
		select {
		case <-s.shutdownCh:
			return
		case <-electionTimer:
			s.SetState(Candidate)
		case rpc := <-s.rpcCh: //把AppendEntriesRequest和VoteRequest放到一个等待队列里，串行执行，防止中间状态错乱
			switch data := rpc.Command.(type) {
			case NoopCommand:
				slog.Warn("follower receive command")
				rpc.Respond(nil, ErrNotLeader)
			case *AppendEntriesRequest:
				electionTimer = randomTimeout(ElectionTimeout) //重置计时器
				resp := s.processAppendEntriesRequest(data)
				rpc.Respond(resp, nil)
			case *VoteRequest:
				resp := s.processVoteRequest(data)
				if resp.Granted { //必须在给对方投票的前提下，才能重置ElectionTimeout计时器
					electionTimer = randomTimeout(ElectionTimeout)
				}
				rpc.Respond(resp, nil)
			case *AppendEntriesResponse:
			case *VoteResponse:
			}
		}
	}
}
func (s *Server) CandidateLoop() {
	slog.Info("server run as candidate", "id", s.Id, "leader id", s.leaderId)
	var leaderChangeTimer <-chan time.Time
	doVote := true
	voteGranted := 0 //获得的票数

	for s.state == Candidate {
		if doVote {
			s.term++          // Term加1
			s.votedFor = s.Id // 投票给自己
			voteGranted++
			lastLogIndex, lastLogTerm := s.log.LastInfo()
			req := &VoteRequest{
				Term:         s.term,
				CandidateId:  s.Id,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			// 异步给每一个peer发送投票请求
			for _, peer := range s.peers {
				s.routineGroup.Add(1)
				go func(peer *Peer) {
					defer s.routineGroup.Done()
					resp, err := s.transporter.RequestVote(peer, req)
					if err == nil {
						rpc := RPC{
							Command: resp,
						}
						s.rpcCh <- rpc //把RequestVote的响应结果放到rpcCh的Command里去
					}
				}(peer)
			}
			doVote = false
			leaderChangeTimer = randomTimeout(LeaderChangeTimeout) //倒计时开始
		}

		if voteGranted >= s.QuorumSize() {
			slog.Info("candidate upgrade to leader", "server id", s.Id, "vote count", voteGranted, "cluster size", len(s.peers)+1, "term", s.term)
			s.SetState(Leader) //成为leader
			return             //退出candidate循环
		}

		select {
		case <-s.shutdownCh:
			return
		case <-leaderChangeTimer:
			slog.Info("leader change timeout")
			doVote = true //如果超时了，则把doVote置为true，把选举流程从头到尾再走一遍
		case rpc := <-s.rpcCh: //把AppendEntriesRequest和VoteRequest放到一个等待队列里，串行执行，防止中间状态错乱
			switch data := rpc.Command.(type) {
			case NoopCommand:
				slog.Warn("candidate receive command")
				rpc.Respond(nil, ErrNotLeader)
			case *AppendEntriesRequest:
				resp := s.processAppendEntriesRequest(data)
				rpc.Respond(resp, nil)
			case *VoteResponse:
				if data.Term > s.term { //对于任何RPC请求或响应，只要对方发过来的Term比自己的大，就无条件地用对方的Term覆盖自己的Term，并把自己降为Follower
					s.upgradeTerm(data.Term, "") //升级term，把votedFor清空，把自己降为follower
					return
				} else {
					if data.Granted {
						voteGranted++
					}
				}
			case *VoteRequest: //也可能会收到其他candidate的投票请求
				resp := s.processVoteRequest(data)
				rpc.Respond(resp, nil)
			case *AppendEntriesResponse:
			}
		}
	}
}

func (s *Server) doHeartBeat() {
	slog.Debug("broadcast heartbeat")
	for _, peer := range s.peers {
		s.routineGroup.Add(1)
		go func(peer *Peer) {
			defer s.routineGroup.Done()
			prevLogIndex := s.prevLogIndex[peer.Id]
			entries, prevTerm := s.log.GetEntriesAfter(prevLogIndex)
			slog.Debug("getEntriesAfter", "follower", peer.Id, "prevLogIndex", prevLogIndex, "log count", len(entries), "prevTerm", prevTerm)
			req := &AppendEntriesRequest{
				Term:         s.term,
				LeaderId:     s.Id,
				CommitIndex:  s.log.CommitIndex(),
				PrevLogIndex: prevLogIndex,
				PrevLogTerm:  prevTerm,
				LogEntries:   entries,
			}
			resp, err := s.transporter.AppendEntries(peer, req)
			if err == nil {
				rpc := RPC{
					Command: resp,
				}
				s.rpcCh <- rpc //把AppendEntries的响应结果放到rpcCh的Command里去
			}
		}(peer)
	}
}

func (s *Server) LeaderLoop() {
	slog.Info("server run as leader", "id", s.Id, "leader id", s.leaderId)
	heartbeatTicker := time.NewTicker(HeartBeatInterval)
	//新leader最开始认为所有follower的prevLogIndex都是自己的lastLogIndex，经过一轮AE（AppendEntries）后，把prevLogIndex置为follower的LastIndex
	lastLogIndex := s.log.LastIndex()
	for _, peer := range s.peers {
		s.prevLogIndex[peer.Id] = lastLogIndex
	}
	s.doHeartBeat() //成为leader后，立即发送一个心跳，让其他candidate放弃。本轮心跳的LogEntries为空

	for s.state == Leader {
		select {
		case <-s.shutdownCh:
			return
		case <-heartbeatTicker.C:
			s.doHeartBeat()
		case rpc := <-s.rpcCh: //把AppendEntriesRequest和VoteRequest放到一个等待队列里，串行执行，防止中间状态错乱
			switch data := rpc.Command.(type) {
			case NoopCommand:
				s.processCommand(data)
			case *AppendEntriesRequest:
				resp := s.processAppendEntriesRequest(data)
				rpc.Respond(resp, nil)
			case *AppendEntriesResponse:
				s.processAppendEntriesResponse(data)
			case *VoteRequest:
				resp := s.processVoteRequest(data)
				rpc.Respond(resp, nil)
			case *VoteResponse:
			default:
				slog.Warn("rpc.NoopCommand", "type", reflect.TypeOf(data).Name())
			}
		}
	}
}

func (s *Server) processAppendEntriesRequest(req *AppendEntriesRequest) *AppendEntriesResponse {
	slog.Debug("AE request", "commit index", req.CommitIndex, "prevIndex", req.PrevLogIndex, "prevTerm", req.PrevLogTerm, "log count", len(req.LogEntries))
	// term比自己的小，拒绝复制日志
	if req.Term < s.term {
		slog.Debug("small term", "leader term", req.Term, "follower term", s.term)
		return &AppendEntriesResponse{Term: s.term, Success: false, LastIndex: s.log.LastIndex(), CommitIndex: s.log.CommitIndex(), FollowerId: s.Id}
	}

	// 多个candidate几乎同时发起投票，其中一个率先成为leader后要立即发一条AppendEntriesRequest，不要等到心跳才发，这样其他candidate收到后立即把自己降为Follower
	if s.state == Candidate {
		slog.Info("downgrade candidate to follower", "server id", s.Id)
	}
	s.upgradeTerm(req.Term, req.LeaderId) // 升级term，把votedFor清空，重置leaderId，把自己降为follower

	if len(req.LogEntries) == 0 {
		return &AppendEntriesResponse{Term: s.term, Success: false, LastIndex: s.log.LastIndex(), CommitIndex: s.log.CommitIndex(), FollowerId: s.Id}
	}

	// 追加日志
	appendSuccess := s.log.AppendEntries(req.PrevLogIndex, req.PrevLogTerm, req.LogEntries)
	// 设置Follower的commit index。不管日志有没有追加成功，都应该去重置commitIndex
	commitIndex := req.CommitIndex       //这是leader在发送AppendEntriesRequest时的commit index
	if s.log.LastIndex() < commitIndex { //由于commitIndex依赖s.log.LastIndex()，所以追加日志要先做，计算commitIndex要后做
		commitIndex = s.log.LastIndex()
	}
	s.log.SetCommitIndex(commitIndex)

	if !appendSuccess {
		return &AppendEntriesResponse{Term: s.term, Success: false, LastIndex: s.log.LastIndex(), CommitIndex: s.log.CommitIndex(), FollowerId: s.Id}
	}
	return &AppendEntriesResponse{Term: s.term, Success: true, LastIndex: s.log.LastIndex(), CommitIndex: s.log.CommitIndex(), FollowerId: s.Id}
}

func (leader *Server) processAppendEntriesResponse(resp *AppendEntriesResponse) {
	slog.Debug("AppendEntriesResponse", "follower id", resp.FollowerId, "Success", resp.Success, "log index", resp.LastIndex, "term", resp.Term, "commit index", resp.CommitIndex)
	// 本来想给follower发送日志，结果发现follower的term比自己的还大，则把自己降为Follower，同时升级term
	if resp.Term > leader.term {
		leader.upgradeTerm(resp.Term, "") // 升级term，把votedFor清空，重置leaderId，把自己降为follower。此时还不能确定谁是Leader，先把LeaderId置空，收到心跳后就知道谁是Leader了。目前代码里也没有使用leaderId
		return
	}

	// 更新follower已同步到了哪条日志。新leader最开始认为所有follower的prevLogIndex都是自己的lastLogIndex，经过一轮AE（AppendEntries）后，把prevLogIndex置为follower的LastIndex
	leader.prevLogIndex[resp.FollowerId] = resp.LastIndex

	if !resp.Success { //follower没有接收AE请求
		return
	}

	//每收到一个Follower的AppendEntriesResponse，Leader都要判断一下是否可以commit一个更大的Index
	var indices []int64                               //存储群集中每个节点的LastLogIndex
	indices = append(indices, leader.log.LastIndex()) //先把Leader的LastLogIndex放进去
	for _, follower := range leader.peers {           //把各个Follower的LastLogIndex放进去
		indices = append(indices, leader.prevLogIndex[follower.Id])
	}
	//所有节点的LastLogIndex按升序排列
	sort.Slice(indices, func(i, j int) bool {
		return indices[i] < indices[j]
	})
	//超过一半的节点都记录到了commitIndex
	commitIndex := indices[leader.QuorumSize()-1]
	slog.Debug("commitIndex", "indices", indices, "commitIndex", commitIndex)
	//新的commitIndex可能不会超过原先的commitIndex
	leader.log.SetCommitIndex(commitIndex)
}

func (s *Server) processVoteRequest(req *VoteRequest) *VoteResponse {
	if req.Term == s.term && len(s.votedFor) > 0 && s.votedFor != req.CandidateId { //同一个term内只能投一次。比如同时出现多个Candidate，它们的Term相同，已经给自己投票了，就不能再投给别人了
		return &VoteResponse{Term: s.term, Granted: false}
	}

	// candidate的term更小，不给它投票
	if req.Term < s.term {
		return &VoteResponse{Term: s.term, Granted: false}
	}
	if req.Term > s.term {
		s.upgradeTerm(req.Term, "") //升级term，把votedFor清空，把自己降为follower
	}

	lastLogIndex, lastLogTerm := s.log.LastInfo()
	// candidate的lastLogIndex/lastLogTerm比follower的还小，则不给candidate投票
	if req.LastLogIndex < lastLogIndex || req.LastLogTerm < lastLogTerm {
		return &VoteResponse{Term: s.term, Granted: false}
	} else { //否则给它投票
		s.upgradeTerm(req.Term, req.CandidateId) //升级term，把votedFor清空，重置leaderId，把自己降为follower
		s.votedFor = req.CandidateId
		slog.Debug("vote grant", "server id", s.Id, "candidate id", req.CandidateId, "term", s.term)
		return &VoteResponse{Term: s.term, Granted: true}
	}
}

func (leader *Server) Do(command NoopCommand) {
	rpc := RPC{
		Command:  command,
		RespChan: make(chan RPCResponse),
	}
	leader.rpcCh <- rpc //把请求放进去
}

// 只有leader才能processCommand
func (leader *Server) processCommand(command NoopCommand) {
	leader.log.CreateEntry(command)
}
