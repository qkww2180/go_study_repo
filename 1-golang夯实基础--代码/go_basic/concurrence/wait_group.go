package concurrence

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wt = sync.WaitGroup{}

func grandson() {
	defer wt.Done()
	fmt.Println("grandson begin")
	fmt.Printf("routine num %d\n", runtime.NumGoroutine())
	time.Sleep(3 * time.Second)
	fmt.Printf("routine num %d\n", runtime.NumGoroutine())
	fmt.Println("grandson finish")
}

func child() {
	defer wt.Done()
	fmt.Println("child begin")
	go grandson()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("child finish") //子协程退出后，孙协挰还在运行。所有协程都是平等的、独立的。一个协程的生命周期不受制于另一个协程（main协程除外）
}

func SubRoutine() { //main协程调用此函数
	wt.Add(2)
	go child()
	// time.Sleep(2 * time.Second)
	wt.Wait()
}

func WaitGroup() {
	const N = 10
	wg := sync.WaitGroup{}
	wg.Add(N) //加N
	for i := 0; i < N; i++ {
		go func(a, b int) { //开N个子协程
			defer wg.Done() //减1
			time.Sleep(10 * time.Millisecond)
			_ = a + b
			fmt.Printf("%d 结束\n", i)
		}(i, i+1)
	}
	fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine()) //N+1
	wg.Wait()                                        //等待减为0
	fmt.Printf("当前协程数：%d\n", runtime.NumGoroutine()) //1
}
