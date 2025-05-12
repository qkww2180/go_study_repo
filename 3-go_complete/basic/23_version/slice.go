package main

import (
	"fmt"
	"slices"
)

func slice() {
	arr := []int{1, 2, 3}
	brr := slices.Repeat(arr, 3)
	fmt.Println(brr) //[1 2 3 1 2 3 1 2 3]
	slices.Sort(brr)
	fmt.Println(brr) //[1 1 1 2 2 2 3 3 3]
	slices.Reverse(brr)
	fmt.Println(brr) //[3 3 3 2 2 2 1 1 1]
}
