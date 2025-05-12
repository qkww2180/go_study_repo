package main

import "fmt"

// 修改一个数字
func modifyInt(a *int) {
	(*a)++
}

// 修改首元素
func modifyEle(s []int) {
	s[0] = 9
}

// 向尾部添加一个元素
func appendEle(s *[]int) {
	*s = append(*s, 9)
}

// 删除尾部元素
func removeEle(s *[]int) {
	n := len(*s)
	*s = (*s)[0 : n-1]
}

func main17() {
	// a := 0
	// modifyInt(&a)
	// fmt.Println(a)

	// s := make([]int, 3)

	// modifyEle(s)
	// fmt.Println(s)

	// appendEle(&s)
	// fmt.Println(s)

	// removeEle(&s)
	// fmt.Println(s)

	arr := [...]int{1, 2, 3, 4, 5}
	brr := arr[1:3:5]
	fmt.Printf("%v %d %d\n", brr, len(brr), cap(brr)) //[2 3] 2 4
}
