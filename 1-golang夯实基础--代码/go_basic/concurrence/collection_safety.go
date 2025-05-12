package concurrence

import (
	"fmt"
	"sync"
)

/**
数组、slice、struct允许并发修改（可能会脏写），并发修改map可能会（不是一定会）发生fatal error，recover()只能捕获panic，但不能捕获fatal error。
如果需要并发修改map请使用sync.Map
*/

type Student struct {
	Name string
	Age  int32
}

var arr = [10]int{}
var m = sync.Map{}
var mp = make(map[int]bool, 10)

func CollectionSafety() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() { //写偶数位
		defer wg.Done()
		for i := 0; i < len(arr); i += 2 {
			arr[i] = 0
		}
	}()
	go func() { //写奇数位
		defer wg.Done()
		for i := 1; i < len(arr); i += 2 {
			arr[i] = 1
		}
	}()
	wg.Wait()
	fmt.Println(arr) //输出[0 1 0 1 0 1 0 1 0 1]
	fmt.Println("=======================")

	wg.Add(2)
	var stu Student
	go func() {
		defer wg.Done()
		stu.Name = "Fred"
	}()
	go func() {
		defer wg.Done()
		stu.Age = 20
	}()
	wg.Wait()
	fmt.Printf("%s %d\n", stu.Name, stu.Age)
	fmt.Println("=======================")

	// fatal error: concurrent map writes
	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	for i := 0; i < 10000; i += 2 {
	// 		mp[i] = true
	// 	}
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	for i := 1; i < 10000; i += 2 {
	// 		mp[i] = false
	// 	}
	// }()
	// wg.Wait()
	// fmt.Printf("%t\n", mp[1])
	// fmt.Println("=======================")

	wg.Add(2)
	go func() {
		defer wg.Done()
		m.Store("k1", "v1")
	}()
	go func() {
		defer wg.Done()
		m.Store("k1", "v2")
	}()
	wg.Wait()
	if v, exists := m.Load("k1"); exists {
		fmt.Println(v)
	}
}
