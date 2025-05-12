package main

import (
	"fmt"
	"slices"
)

func main31() {
	arr := []int{3, 7, 4, 6, 8, 1}
	slices.Sort(arr)
	fmt.Println(arr)
	slices.SortFunc(arr, func(a, b int) int { //自定义排序方式
		return b - a
	})
	fmt.Println(arr)

	type User struct {
		Age    int
		Height float32
	}
	brr := []*User{&User{18, 1.8}, &User{25, 1.7}}
	slices.SortFunc(brr, func(a, b *User) int {
		// return int(b.Height - a.Height)
		if b.Height > a.Height {
			return 1
		} else if b.Height < a.Height {
			return -1
		} else {
			return 0
		}
	})

	fmt.Println("最大者", slices.Max(arr))
	fmt.Println("最小者", slices.Min(arr))
	fmt.Println("包含", slices.Contains(arr, 5))

	crr := make([]int, len(arr))
	copy(crr, arr) //最多只能拷贝 len(crr) 个元素，性能比自己写for循环要高很多
	fmt.Println(crr)

	fmt.Println("相等", slices.Equal(arr, crr)) //true
	arr[0]++
	fmt.Println("相等", slices.Equal(arr, crr)) //false

	drr := arr
	fmt.Println("相等", slices.Equal(drr, arr)) //true
	arr[0]++
	fmt.Println("相等", slices.Equal(drr, arr)) //true
}
