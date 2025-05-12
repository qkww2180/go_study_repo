package main

import "fmt"

var g int = 7

type x int //x是int的别名吗？

func main1() {
	f := 1. //不能有隐式的数字类型转换
	p := &g //获取全局变量的地址
	//不允许对指针施加算术运算
	fmt.Printf("%.2f %p %d %p\n", f, p, *p, p)
	g--  // ++、--只是赋值操作，不是表达式
	*p-- //修改栈变量的值，这不是对指针施加运算

	var a int
	var b x //a和b是相同的数据类型吗？
	fmt.Println(a, b)

	b.dance()
	// a.dance()
}

func (x) dance() {}
