package main

import (
	"fmt"
)

/**
何时会发生panic:
1. index out of range和divide by zero是常见的runtime error,会发生panic
2. 通常系统初始化发生重大问题会主动调用panic(any)

panic时会依次执行：
1. 执行已经注册的defer(后注册的先执行)，未注册的defer不执行
2. 打印错误信息和堆栈调用信息
3. 调用os.Exit(2)结束go进程
*/

func defer_panic() {
	defer fmt.Println(1)
	var arr []int
	n := 0
	// defer fmt.Println(1 / n) //在注册defer时就要计算1/n，发生panic
	defer func() {
		_ = arr[n]
		_ = 1 / n            //defer func 内部发生panic，main协程不会exit，其他defer还可以正常执行
		defer fmt.Println(2) //上面那行代码发生发panic，所以本行的defer没有注册成功
	}()
	defer fmt.Println(3)
}

// recover()一定要放在函数的最开始位置
func soo() {
	// defer recover()//recover()必须放在defer func() {}里,不能直接放在defer后面,否则不生效
	defer func() { //去掉这个defer试试，看看panic的流程。把这个defer放到soo函数末尾试试。把这个defer移到main()里试试。
		//recover必须在defer中才能生效
		if panicInfo := recover(); panicInfo != nil {
			fmt.Printf("soo函数中发生了panic:%v\n", panicInfo)
			// debug.PrintStack() //打印调用堆栈
		}
	}()

	fmt.Println("enter soo")

	fmt.Println("regist recover")

	defer fmt.Println("hello")
	defer func() {
		n := 0
		_ = 3 / n //除0异常，发生panic，下一行的defer没有注册成功
		defer fmt.Println("how are you")
	}()
}

func main17() {
	// defer_panic()
	soo()
}
