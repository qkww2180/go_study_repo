package concurrence

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var n int32 = 0

func inc1() {
	n++ //n++不是原子操作，它分为3步：取出n，加1，结果赋给n
}

func inc2() {
	atomic.AddInt32(&n, 1) //封装成原子操作
}

func Atomic() {
	const P = 1000 //开大量协程才能把脏写问题测出来
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc1()
		}()
	}
	wg.Wait()
	fmt.Printf("finally n=%d\n", n) //多运行几次，n经常不等于1000
	fmt.Println("===========================")
	n = 0 //重置n
	wg = sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc2()
			_ = n + 1
			_ = atomic.LoadInt32(&n) + 1 //即使在高并发环境下，读一个变量也没必要用atomic.Load
		}()
	}
	wg.Wait()
	fmt.Printf("finally n=%d\n", n)
}
