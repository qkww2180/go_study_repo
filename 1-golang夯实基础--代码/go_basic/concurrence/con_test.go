package concurrence_test

import (
	"dqq/go/basic/concurrence"
	"sync"
	"testing"
)

func TestSimpleGoroutine(t *testing.T) {
	concurrence.SimpleGoroutine()
}

func TestWaitGroup(t *testing.T) {
	concurrence.SubRoutine()
	concurrence.WaitGroup()
}

func TestLock(t *testing.T) {
	// concurrence.Atomic()
	// concurrence.RWlock()
	// concurrence.ReentranceRLock(3)
	// concurrence.ReentranceWLock(3)
	concurrence.RLockExclusion()
	// concurrence.WLockExclusion()
}

func TestCollectionSafety(t *testing.T) {
	concurrence.CollectionSafety()
}

func TestConcurrentMap(t *testing.T) {
	cm1 := concurrence.NewConcurrentMap[string, int](50)
	cm1.Store("张三", 18)
	if v, exists := cm1.Load("张三"); !exists {
		t.Fail()
	} else {
		if v != 18 {
			t.Fail()
		}
	}
	if _, exists := cm1.Load("李四"); exists {
		t.Fail()
	}

	cm2 := concurrence.NewConcurrentMap[int, bool](50)
	cm2.Store(18, true)
	if v, exists := cm2.Load(18); !exists {
		t.Fail()
	} else {
		if v != true {
			t.Fail()
		}
	}
	if _, exists := cm2.Load(19); exists {
		t.Fail()
	}

	// 测试在高并发情况下是否安全
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				cm2.Store(j, true)
				cm2.Load(j)
			}
		}()
	}
	wg.Wait()
}

// go test -v ./concurrence -run=^TestSimpleGoroutine$ -count=1
// go test -v ./concurrence -run=^TestWaitGroup$ -count=1
// go test -v ./concurrence -run=^TestLock$ -count=1
// go test -v ./concurrence -run=^TestCollectionSafety$ -count=1
// go test -v ./concurrence -run=^TestConcurrentMap$ -count=1
