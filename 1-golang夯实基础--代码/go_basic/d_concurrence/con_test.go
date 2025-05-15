package d_concurrence_test

import (
	"dqq/go/basic/d_concurrence"
	"sync"
	"testing"
)

func TestSimpleGoroutine(t *testing.T) {
	d_concurrence.SimpleGoroutine()
}

func TestWaitGroup(t *testing.T) {
	d_concurrence.SubRoutine()
	d_concurrence.WaitGroup()
}

func TestLock(t *testing.T) {
	// d_concurrence.Atomic()
	// d_concurrence.RWlock()
	// d_concurrence.ReentranceRLock(3)
	// d_concurrence.ReentranceWLock(3)
	d_concurrence.RLockExclusion()
	// d_concurrence.WLockExclusion()
}

func TestCollectionSafety(t *testing.T) {
	d_concurrence.CollectionSafety()
}

func TestConcurrentMap(t *testing.T) {
	cm1 := d_concurrence.NewConcurrentMap[string, int](50)
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

	cm2 := d_concurrence.NewConcurrentMap[int, bool](50)
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

// go test -v ./d_concurrence -run=^TestSimpleGoroutine$ -count=1
// go test -v ./d_concurrence -run=^TestWaitGroup$ -count=1
// go test -v ./d_concurrence -run=^TestLock$ -count=1
// go test -v ./d_concurrence -run=^TestCollectionSafety$ -count=1
// go test -v ./d_concurrence -run=^TestConcurrentMap$ -count=1
