package raft

import (
	"log/slog"
	"sync"
)

type LogEntry struct {
	Term        int64
	Index       int64
	NoopCommand NoopCommand
}

type Log struct {
	sync.RWMutex
	server  *Server
	entries []*LogEntry //LogEntry的Term和Index都是升序的
	// file        os.File     //日志写内存的同时就要写磁盘，宕机重启后可以从磁盘恢复日志
	startIndex  int64 //从1开始，0表示日志集合为空
	commitIndex int64 //初始为0
}

func NewLog(s *Server) *Log {
	return &Log{entries: make([]*LogEntry, 0, 100), server: s, startIndex: 1}
}

// 只有在leader上能调用CreateEntry
func (log *Log) CreateEntry(command NoopCommand) *LogEntry {
	log.Lock()
	defer log.Unlock()
	var lastIndex int64
	if len(log.entries) > 0 {
		entry := log.entries[len(log.entries)-1]
		lastIndex = entry.Index
	}
	entry := &LogEntry{Term: log.server.term, Index: lastIndex + 1, NoopCommand: command}
	log.entries = append(log.entries, entry)
	slog.Debug("append log", "index", entry.Index, "term", entry.Term)
	return entry
}

func (log *Log) LastIndex() int64 {
	log.RLock()
	defer log.RUnlock()
	if len(log.entries) == 0 {
		return 0
	}
	return log.entries[len(log.entries)-1].Index
}

func (log *Log) LastInfo() (int64, int64) {
	log.RLock()
	defer log.RUnlock()
	if len(log.entries) == 0 {
		return 0, 0
	}
	entry := log.entries[len(log.entries)-1]
	return entry.Index, entry.Term
}

func (log *Log) CommitIndex() int64 {
	log.RLock()
	defer log.RUnlock()
	return log.commitIndex
}

// 根据Index找对应的LogEntry
func (log *Log) findByIndex(target int64) (int, *LogEntry) {
	// log index可能是不连续的，但肯定是递增的，可以用二分查找法
	arr := log.entries
	begin := 0
	end := len(arr) - 1
	for begin <= end { //之所以是<=，而不是<，是因为区间内只剩下一个元素时也应该跟target进行比较，而不是直接返回-1
		middle := (begin + end) / 2
		if arr[middle].Index == target {
			return middle, arr[middle]
		}
		if arr[middle].Index < target {
			begin = middle + 1
		} else {
			end = middle - 1
		}
	}
	slog.Debug("could not found log by index", "idx", target)
	return -1, nil

	// log index就是在数组里的下标，是连续的，没必要用二分查找
	// target -= log.startIndex
	// if target >= 0 && target < int64(len(log.entries)) {
	// 	return int(target), log.entries[target]
	// } else {
	// 	return -1, nil
	// }
}

// 提交第idx条日志，及其之前的日志
func (log *Log) SetCommitIndex(idx int64) int64 {
	log.Lock()
	defer log.Unlock()

	i, entry := log.findByIndex(idx)
	if entry == nil {
		slog.Warn("could not found log entry", "log index", idx)
		return -1
	}

	// 只能提交本Term的index
	if log.server.state == Leader && log.server.term != entry.Term {
		return -1
	}

	// commit index只能比之前的大
	if idx == log.commitIndex {
		return -1
	}
	if idx < log.commitIndex {
		slog.Warn("current commit index less than prev commit index", "current commit index", idx, "prev commit index", log.commitIndex)
		return -1
	}

	prevCommitIndex := log.commitIndex
	//更新commit index
	log.commitIndex = idx
	// 从上一次的commit index到本次的commit index，这之间的log的command应用到状态机
	j, _ := log.findByIndex(prevCommitIndex)
	if j < 0 {
		return -1
	}
	slog.Debug("apply log", "from", prevCommitIndex, "to", idx)
	for k := j + 1; k <= i; k++ {
		entry := log.entries[k]
		entry.NoopCommand.Apply(log.server.fsm)
	}
	return idx
}

// 从自己的log集合里找到<prevLogIndex, prevLogTerm>这条日志，然后把entries追加到它后面。<prevLogIndex, prevLogTerm>没有被覆盖
func (log *Log) AppendEntries(prevLogIndex, prevLogTerm int64, entries []*LogEntry) bool {
	log.Lock()
	defer log.Unlock()
	if len(entries) == 0 {
		return false
	}
	total := len(log.entries)
	if prevLogIndex == 0 {
		log.entries = entries
		return true
	}

	i, entry := log.findByIndex(prevLogIndex)
	if i < 0 {
		return false
	}
	//commit过的日志不能被覆盖
	if entry.Index < log.startIndex || entry.Index < log.commitIndex {
		return false
	}
	// term和index可以唯一确定一条日志
	if entry.Term == prevLogTerm {
		if i < total-1 {
			log.entries = log.entries[:i+1] //丢弃后半部分(第i条不丢弃)
		}
		log.entries = append(log.entries, entries...) //把新日志追加到后面
		slog.Debug("AppendEntries", "from", prevLogIndex, "count", len(entries))
		return true
	} else { //index相同，但term不同
		log.entries = log.entries[:i] //丢弃后半部分(包括第i条也丢弃)
		slog.Debug("AppendEntries could not found match", "prevLogIndex", prevLogIndex, "prevLogTerm", prevLogTerm)
		return false
	}
}

// 获取prevLogIndex往后的日志，不包含prevLogIndex
func (log *Log) GetEntriesAfter(prevLogIndex int64) ([]*LogEntry, int64) {
	log.Lock()
	defer log.Unlock()
	if len(log.entries) == 0 {
		return nil, 0
	}
	if prevLogIndex >= log.entries[len(log.entries)-1].Index {
		return nil, 0
	}
	total := len(log.entries)
	var entries []*LogEntry
	if prevLogIndex < log.startIndex {
		entries = log.entries
	} else {
		i, _ := log.findByIndex(prevLogIndex)
		if i < total-1 {
			entries = log.entries[i+1:]
		}
	}

	if len(entries) > MaxLogEntriesPerHeartbeat {
		entries = entries[:MaxLogEntriesPerHeartbeat]
	}

	if len(entries) == 0 {
		return nil, 0
	} else {
		return entries, entries[len(entries)-1].Term
	}
}
