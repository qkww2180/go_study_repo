package raft

import "errors"

// Finite State Machine
type FSM interface {
}

// no operation。一个空的command
type NoopCommand struct {
}

func (NoopCommand) Apply(sm FSM) error {
	return nil
}

type VoteRequest struct {
	Term         int64
	CandidateId  string
	LastLogIndex int64
	LastLogTerm  int64
}

type VoteResponse struct {
	Term    int64
	Granted bool
}

type AppendEntriesRequest struct {
	Term         int64
	LeaderId     string
	CommitIndex  int64
	PrevLogIndex int64
	PrevLogTerm  int64
	LogEntries   []*LogEntry
}

type AppendEntriesResponse struct {
	Term        int64
	Success     bool
	LastIndex   int64
	CommitIndex int64 //没用上，只上为了打日志
	FollowerId  string
}

type State uint32

const (
	Follower State = iota
	Candidate
	Leader
	Stopped
)

// go语言打印变量时会默认调用变量的String()方法
func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

var (
	ErrNotLeader    = errors.New("not leader")
	ErrNotCandidate = errors.New("not candidate")
	ErrShutdown     = errors.New("has been shutdown")
)
