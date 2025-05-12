package main

import (
	"sync"
)

var (
	pool = sync.Pool{
		New: func() any {
			return make([]int, 10000)
		},
	}
)

func Search(keyword string) []int {
	result := pool.Get().([]int)
	for i := 0; i < 10000; i++ {
		result[i] = i
	}
	return result
}
