package test

import (
	"dqq/algorithm"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAliasSample(t *testing.T) {
	endpoints := []string{"127.0.0.1:1234", "127.0.0.1:2345", "127.0.0.1:3456"}
	probs := []float64{1, 2, 3} //每个endpoint被抽中的概率
	useCount := make([]int32, len(endpoints))
	balancer := algorithm.NewAliasSampler(probs)
	const P = 100 //开100个协程并发使用balancer
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				idx := balancer.Sample() //取出一个endpoint
				atomic.AddInt32(&useCount[idx], 1)
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond) //假装在使用endpoint
			}
		}()
	}
	wg.Wait()

	fmt.Println(useCount) //打印每个endpoint被使用了几次
}

// go test ./algorithm/test -v -run=^TestAliasSample$ -count=1
