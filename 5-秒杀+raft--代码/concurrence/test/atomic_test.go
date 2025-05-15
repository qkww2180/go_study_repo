package test

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
)

func TestReadDirect(t *testing.T) {
	const L int = 10
	arr := [L]uint64{}
	var i uint64
	for i = 0; i < 100; i++ {
		wg := sync.WaitGroup{}
		wg.Add(L + 1)
		var v uint64
		go func() { defer wg.Done(); v = math.MaxUint64 - i }()
		for j := 0; j < L; j++ {
			go func(j int) { defer wg.Done(); arr[j] = v }(j)
		}
		wg.Wait()

		for j := 0; j < L; j++ {
			if arr[j] != 0 && arr[j] != math.MaxUint64-i {
				fmt.Println(arr[j])
				t.Fail()
				break
			}
		}
	}
}

func TestAtomicLoad(t *testing.T) {
	const L int = 10
	arr := [L]uint64{}
	var i uint64
	for i = 0; i < 100; i++ {
		wg := sync.WaitGroup{}
		wg.Add(L + 1)
		var v uint64
		go func() { defer wg.Done(); atomic.StoreUint64(&v, math.MaxUint64-i) }()
		for j := 0; j < L; j++ {
			go func(j int) { defer wg.Done(); arr[j] = atomic.LoadUint64(&v) }(j)
		}
		wg.Wait()

		for j := 0; j < L; j++ {
			if arr[j] != 0 && arr[j] != math.MaxUint64-i {
				fmt.Println(arr[j])
				t.Fail()
				break
			}
		}
	}
}

// go test -v ./d_concurrence/test -run=^TestReadDirect$ -count=1
// go test -v ./d_concurrence/test -run=^TestAtomicLoad$ -count=1
