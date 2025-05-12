package test

import (
	"context"
	"dqq/algorithm"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

var (
	totalQuery int32
	limiter    = rate.NewLimiter(rate.Every(100*time.Millisecond), 1) //每隔100ms生成一个令牌，桶的容量为1，最大QPS限制为10
	myLimiter  = algorithm.NewLimiter(10, 1)                          //1秒生成10个令牌，桶的容量为1，最大QPS限制为10
)

func handler() {
	atomic.AddInt32(&totalQuery, 1)
	time.Sleep(50 * time.Millisecond)
}

func wait() {
	limiter.WaitN(context.Background(), 1) //阻塞，直到桶中有N个令牌。N=1时等价于Wait(context)
	handler()
}

func myWait() {
	myLimiter.WaitN(1) //阻塞，直到桶中有N个令牌。N=1时等价于Wait(context)
	handler()
}

func allow() {
	if limiter.AllowN(time.Now(), 1) { //当前桶中是否至少还有N个令牌，如果有则返回true，否则返回false。N=1时等价于Allow(time.Time)
		handler()
	}
}

func reserve() {
	reserve := limiter.ReserveN(time.Now(), 1)
	time.Sleep(reserve.Delay()) //reserve.Delay()告诉你还需要等多久才会有充足的令牌，你就等吧
	handler()
}

func printQPS() {
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Printf("过去1秒钟接口调用了%d次\n", atomic.LoadInt32(&totalQuery))
		atomic.StoreInt32(&totalQuery, 0) //每隔一秒，清0一次
	}
}

func TestWait(t *testing.T) {
	go func() {
		for {
			wait()
		}
	}()
	printQPS()
}

func TestAllow(t *testing.T) {
	go func() {
		for {
			allow()
		}
	}()
	printQPS()
}

func TestReserve(t *testing.T) {
	go func() {
		for {
			reserve()
		}
	}()
	printQPS()
}

func TestMyWait(t *testing.T) {
	go func() {
		for {
			wait()
		}
	}()
	printQPS()
}

// go test ./algorithm/test -v -run=^TestWait$ -count=1
// go test ./algorithm/test -v -run=^TestAllow$ -count=1
// go test ./algorithm/test -v -run=^TestReserve$ -count=1
// go test ./algorithm/test -v -run=^TestMyWait$ -count=1
