package grpcconnpool

import "sync"

type LoadBalancer interface {
	Select(n int) int
}

type RoundRobin struct {
	currIndex int
	mu        sync.Mutex
}

func (rr *RoundRobin) Select(n int) int {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	i := rr.currIndex
	rr.currIndex = (i + 1) % n
	return i
}
