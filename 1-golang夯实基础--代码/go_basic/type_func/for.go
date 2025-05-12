package main

import (
	"fmt"
	"math/rand"
	"time"
)

func for_basic() {
	arr := []int{1, 2, 3, 4, 5}
	//for 初始化局部变量;条件表达式;后续操作
	for i := 0; i < len(arr); i++ { //正序遍历切片
		fmt.Printf("%d: %d\n", i, arr[i])
	}
	fmt.Println()
	for i := len(arr) - 1; i >= 0; i-- { //逆序遍历切片
		fmt.Printf("%d: %d\n", i, arr[i])
	}
	fmt.Println()
	//初始化多个变量; 复制的逻辑表达式; 后续操作有多个
	for sum, i := 0, 0; i < len(arr) && sum < 100; sum, i = sum*1, i+1 {
		sum += arr[i] //累积和
		fmt.Printf("%d: %d\n", i, sum)
	}
}

func part_for() {
	arr := []int{1, 2, 3, 4, 5}
	var i int = 0

	//把初始化工作放到for上面
	for ; i < len(arr); i++ {
		fmt.Printf("%d: %d\n", i, arr[i])
	}
	fmt.Println()

	//把后续操作放到for块内部
	for j := 0; j < len(arr); {
		fmt.Printf("%d: %d\n", j, arr[j])
		j++
	}
	fmt.Println()

	//for后面只跟一个条件表达式
	i = 0
	for i < len(arr) { //演示：把分号加上
		fmt.Printf("%d: %d\n", i, arr[i])
		i++
	}
	fmt.Println()

	//for后面只跟一个条件表达式时，分号可以不要
	i = 0
	for i < len(arr) {
		fmt.Printf("%d: %d\n", i, arr[i])
		i++
	}
	fmt.Println()
}

// 无限循环
func endless_loop() {
	for { //等价于for true
		fmt.Print("zzz...")
		time.Sleep(2 * time.Second)
	}
}

func for_range() {
	//数组、切片、string、map、channel都可以通过for range的形式遍历
	arr := [5]int{1, 2, 3, 4, 5}
	for i, ele := range arr {
		fmt.Printf("%d: %d\n", i, ele)
	}
	fmt.Println()

	slc := []int{1, 2, 3, 4, 5}
	for i, ele := range slc {
		fmt.Printf("%d: %d\n", i, ele)
	}
	fmt.Println()

	str := "我会唱ABC"
	for i, ele := range str {
		fmt.Printf("%d: %d %s\n", i, ele, string(ele)) //ele是rune类型，而rune实际上是int32，所以要用%d
	}

	m := map[int]int{1: 2, 3: 6, 5: 10, 7: 14, 9: 18, 11: 22}
	//多次遍历map返回的顺序是不一样的，但相对顺序是一样的，因为每次随机选择一个开始位置，然后顺序遍历
	for key, value := range m {
		fmt.Printf("%d=%d\n", key, value)
	}
	fmt.Println("---------")
	for key, value := range m {
		fmt.Printf("%d=%d\n", key, value)
	}

	ch := make(chan int, 10)
	for i := 4; i < 16; i += 2 {
		ch <- i
	}
	close(ch) //遍历channel之前，一定要先close掉，确保不会再写入
	for ele := range ch {
		fmt.Println(ele)
	}
}

func nest_for() {
	const SIZE = 4

	A := [SIZE][SIZE]float64{}
	//初始化二维数组
	//两层for循环嵌套
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			A[i][j] = rand.Float64() //[0,1)上的随机数
		}
	}

	B := [SIZE][SIZE]float64{}
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			B[i][j] = rand.Float64() //[0,1)上的随机数
		}
	}

	rect := [SIZE][SIZE]float64{}
	//三层for循环嵌套
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			prod := 0.0
			for k := 0; k < SIZE; k++ {
				prod += A[i][k] * B[k][j]
			}
			rect[i][j] = prod
		}
	}

	i, j := 2, 1
	fmt.Println(A[i]) //二维数组第i行
	//打印二维数组的第j列
	//注意：B[:][j]这不是二维数组第j列，这是二维数组第j行！
	for _, row := range B {
		fmt.Printf("%g ", row[j])
	}
	fmt.Println()
	fmt.Println(rect[i][j])
}
