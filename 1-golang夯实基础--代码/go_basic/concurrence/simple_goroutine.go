package concurrence

import (
	"fmt"
	"runtime"
	"time"
)

func Add(a, b int) int {
	fmt.Println("Add")
	return a + b
}

var add = func(a, b int) int {
	fmt.Println("add")
	return a + b
}

func SimpleGoroutine() {
	fmt.Printf("逻辑处理器数目:%d\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU() / 2) //设置本进程最多使用几个核（这里指MPG中的P）

	go Add(2, 5)
	go Add(2, 5)
	go func(a, b int) int {
		fmt.Println("Add")
		return a + b
	}(3, 7)
	go add(4, 8)
	fmt.Printf("存在的协程数目 %d\n", runtime.NumGoroutine())

	time.Sleep(1 * time.Millisecond)
}
