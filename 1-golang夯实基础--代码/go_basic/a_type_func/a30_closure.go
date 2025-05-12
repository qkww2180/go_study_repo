package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Go 函数可以是一个闭包。闭包是一个函数值，它引用了函数体之外的变量。 这个函数可以对这个引用的变量进行访问和赋值；换句话说这个函数被“绑定”在这个变量上。
*/

func f() int {
	a := 3
	a++
	return a
}

func addd() func() int {
	a := 3
	fmt.Printf("in addd func, address of a is %p\n", &a)
	return func() int { //为这个函数分配了一个a。这个函数持有了a。a从栈上逃逸到了堆上
		fmt.Printf("in closure func, address of a is %p\n", &a)
		a++
		return a
	}
}

// 斐波那契序列
func fibonacci() func() int {
	a, b := 0, 1 //通常情况，函数内的局部变量是分配到栈上。但是此处a,b会逃逸到堆上
	return func() int {
		a, b = b, a+b
		_ = b
		return a
	}
}

func fibonacci_seq(n int) []int {
	f := fibonacci() //闭包是一个函数值
	seq := make([]int, 0, n)
	for i := 0; i < n; i++ {
		seq = append(seq, f())
	}
	return seq
}

func funcTimer() {
	begin := time.Now()
	// defer fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())   // defer会立刻计算后面的结果
	defer func() { //defer 后面跟闭包。闭包里引用了begin变量
		fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds()) // 传给defer的是一个函数指针，还没有真正执行这个defer
	}()
	time.Sleep(time.Second)
}

func forRange() {
	arr := []int{1, 2, 3, 4}
	wg := sync.WaitGroup{}
	wg.Add(len(arr))
	for _, v := range arr {
		go func() { //go 后面跟闭包。闭包里引用了v变量（每一轮for循环会创建全新的v）,通常进入下轮循环时上一轮的v变量就不可访问了，但由于闭包，它仍被闭包函数持有
			defer wg.Done()
			time.Sleep(time.Second)
			fmt.Printf("%d\t", v) //输出 1 2 3 4 (顺序随机)
		}()
	}
	wg.Wait()
	fmt.Println()
}

func funcList() {
	var funcSlice []func()
	for i := 0; i < 3; i++ {
		fmt.Printf("in for loop, address of i is %p\n", &i)
		funcSlice = append(funcSlice,
			func() { //该函数持有了i
				fmt.Printf("in closure func, address of i is %p\n", &i)
				println(i)
			},
		)
	}
	for j := 0; j < 3; j++ {
		funcSlice[j]() // 0 1 2
	}
}

func main30() {
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println()
	af := addd()
	fmt.Println(af())
	fmt.Println(af())
	fmt.Println()
	bf := addd()
	fmt.Println(bf())
	fmt.Println(bf())
	fmt.Println()
	funcList()

	fmt.Println(fibonacci_seq(10))
	forRange()
}
