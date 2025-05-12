package main

import (
	"fmt"
	"sync"
)

type Author struct {
	Male bool
}

type Book struct {
	Price   float32
	Name    string
	Author1 Author
	Author2 *Author
}

type MyInterface interface{}

var lock1 sync.Mutex     //go变量声明即初始化，跟C语言不一样
var lock2 = sync.Mutex{} //跟上一行等价

func main14() {
	lock1.Lock()
	defer lock1.Unlock()
	lock2.Lock()
	defer lock2.Unlock()

	var a int
	var b bool
	var c float32
	var d string
	var e any
	var f MyInterface
	var g *int
	fmt.Printf("%T  %[1]v\n", a)
	fmt.Printf("%T  %[1]v\n", b)
	fmt.Printf("%T  %[1]v\n", c)
	fmt.Printf("%T  【%[1]v】\n", d)
	fmt.Printf("%T  %[1]v\n", e)
	fmt.Printf("%T  %[1]v\n", f)
	fmt.Printf("%T  %[1]v\n", g)

	var h Book
	fmt.Printf("%T  %+[1]v\n", h)
}
