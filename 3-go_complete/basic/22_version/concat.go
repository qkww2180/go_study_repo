package main

import (
	"fmt"
	"slices"
)

func main2() {
	arr1 := []int{1, 2, 3}
	arr2 := []int{4, 5, 6}
	arr3 := []int{7, 8, 9}

	merged := append(arr1, arr2...) //v1.22之前合并切片
	merged = append(merged, arr3...)
	fmt.Println(merged)

	merged = slices.Concat(arr1, arr2, arr3) //v1.22更方便，并且是基于泛型的
	fmt.Println(merged)
}
