package main

import (
	"errors"
	"fmt"
	"time"
)

func main10() {
	c, d := 3, 5
	arg1(c, d)
	arg2(c, d)
	arg3(&c, &d)
	return1(c, d)
	return2(c, d)
	return3()
	variable_ength_arg(c, d, 5)
	fmt.Println(c, d)
	fmt.Println(sum(1, 2, 3, 4)) //10
	fmt.Println(sum(1))          //1
	fmt.Println(sum())           //0
}

// a,b是形参，形参是函数内部的局部变量，实参的值会拷贝给形参
func arg1(a int, b int) { //注意大括号{不能另起一行
	a = a + b //在函数内部修改形参的值，实参的值不受影响
	return    //函数返回，return后面的语句不会再执行
	fmt.Println("我不会被输出")
}

func arg2(a, b int) { //参数类型相同时可以只写一次
	a = a + b
	//不写return时，默认执行完最后一行代码函数返回
}

func arg3(a, b *int) { //如果想通过函数修改实参，就需要指针类型
	*a = *a + *b
	*b = 888
}

func no_arg() { //函数可以没有参数，也没有返回值
	fmt.Println("欢迎开启Golang之旅")
}

func return1(a, b int) int { //函数需要返回一个int型数据
	a = a + b
	c := a //声明并初始化一个变量c
	return c
}

func return2(a, b int) (c int) { //返回变量c已经声明好了
	a = a + b
	c = a  //直接使用c
	return //由于函数要求有返回值，即使给c赋过值了，也需要显式写return
}

func return3() (int, int) { //可以没有形参，可以返回多个参数
	now := time.Now()
	return now.Hour(), now.Minute()
}

// 不定长参数
func variable_ength_arg(a int, other ...int) int { //调用该函数时，other可以对应0个参数也可以对应多个参数
	sum := a
	//不定长参数实际上是slice类型
	for _, ele := range other {
		sum += ele
	}
	if len(other) > 0 {
		fmt.Printf("first ele %d len %d cap %d\n", other[0], len(other), cap(other))
	} else {
		fmt.Printf("len %d cap %d\n", len(other), cap(other))
	}
	return sum
}

func sum(arr ...int) int {
	s := 0
	if len(arr) == 0 {
		return s
	}
	s += arr[0]
	s += sum(arr[1:]...)
	return s
}

// 计算斐波那契数列的第n个值
//
// 斐波那契数列以如下被以递推的方法定义： F (0)=0， F (1)=1, F (n)= F (n - 1)+ F (n - 2)
//
// 斐波那契数列前10个数为：0，1，1，2，3，5，8，13，21，34
func Fibonacci(n int) int {
	if n == 0 || n == 1 {
		return n //凡是递归，一定要有终止条件，否则会进入无限循环
	}
	return Fibonacci(n-1) + Fibonacci(n-2) //递归调用函数自身
}

// 递归实现所有参数乘积的倒数。
// error一般作为最后一个返回值
func div(args ...float64) (float64, error) {
	if len(args) == 0 {
		return 0, errors.New("divide by zero")
	}
	first := args[0]
	if first == 0 {
		return 0, errors.New("divide by zero")
	}
	if len(args) == 1 {
		return 1 / first, nil
	}
	remain := args[1:]
	res, err := div(remain...)
	if err != nil {
		return 0, err
	} else {
		return 1 / first * res, nil
	}
}
