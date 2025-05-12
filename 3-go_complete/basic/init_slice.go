package main

import "fmt"

func main15() {
	arr := []int{}
	brr := make([]int, 3, 100)

	arr = append(arr, 1)
	brr = append(brr, 1)

	fmt.Println(arr)
	fmt.Println(brr)
}
