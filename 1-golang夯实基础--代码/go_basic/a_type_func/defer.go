package main

import (
	"fmt"
	"time"
)

// defer典型的应用场景是释放资源，比如关闭文件句柄，释放数据库连接等

func basicDefer() {
	fmt.Println("A")
	defer fmt.Println(1) //defer用于注册一个延迟调用（在函数返回之前调用）
	fmt.Println("B")
	defer fmt.Println(2) //如果同一个函数里有多个defer，则后注册的先执行
	fmt.Println("C")
	defer fmt.Println(3)
	fmt.Println("D")
}

func defer_exe_time() (i int) {
	i = 9
	defer func() { //defer后可以跟一个func
		fmt.Printf("first i=%d\n", i) //打印5，而非9。充分理解“defer在函数返回前执行”的含义，不是在“return语句前执行defer”
	}()
	defer func(i int) {
		fmt.Printf("second i=%d\n", i) //打印9。在注册defer时i就传进去了
	}(i)
	defer fmt.Printf("third i=%d\n", i) //defer后不是跟func，而直接跟一条执行语句，则相关变量在注册defer时被拷贝或计算
	return 5
}

func timeOfWork(arg int) int {
	begin := time.Now()
	defer func() { fmt.Printf("use time %f seconds\n", time.Since(begin).Seconds()) }()
	if arg > 10 {
		time.Sleep(2 * time.Second)
		return 100
	} else {
		time.Sleep(3 * time.Second)
		return 200
	}
}

func main16() {
	basicDefer()
	defer_exe_time()
	timeOfWork(13)
}
