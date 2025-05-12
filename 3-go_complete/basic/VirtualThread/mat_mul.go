package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// 任务本身
func task() {
}

func main() {
	// 不需要热身
	p, _ := strconv.Atoi(os.Args[1]) //p个并发协程
	fmt.Printf("concurrency %d\n", p)
	wg := sync.WaitGroup{}
	wg.Add(p)
	begin := time.Now() //开始计时
	for i := 0; i < p; i++ {
		go func() {
			defer wg.Done()
			task() // 每个协程执行一次任务
		}()
	}
	wg.Wait()                                   //等所有任务结束
	useTime := time.Since(begin).Milliseconds() //结束计时
	fmt.Printf("use time %d ms\n", useTime)
}

// go build -o basic/VirtualThread/matmul.exe basic/VirtualThread/mat_mul.go
// ./basic/VirtualThread/matmul.exe 10000
