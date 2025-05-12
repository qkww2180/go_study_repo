package main

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"maps" //maps从golang.org/x/exp包下移到了标准库里
)

// map按Key排序
func sortMap() {
	m := map[string]struct{}{"赵六": struct{}{}, "张三": struct{}{}, "王五": struct{}{}, "李四": struct{}{}}
	for _, key := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%s\t", key)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
}

func seq() {
	m := map[string]struct{}{"赵六": struct{}{}, "张三": struct{}{}, "王五": struct{}{}, "李四": struct{}{}}
	s1 := maps.Keys(m) //不保证每次的顺序都一样
	// Seq是这种函数：type Seq[V any] func(yield func(V) bool)
	//
	// 可以通过for range 直接遍历这种函数
	for key := range s1 {
		fmt.Println(key)
	}
	fmt.Println(strings.Repeat("-", 50))

	// 也可以借助于Pull()和next()遍历Seq
	next, stop := iter.Pull(s1)
	defer stop()
	for {
		key, valid := next()
		if valid {
			fmt.Println(key)
		} else {
			break
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}

func seq2() {
	m := map[string]string{"赵六": "赵六", "张三": "张三", "王五": "王五", "李四": "李四"}
	s2 := maps.All(m)
	next, stop := iter.Pull2(s2)
	defer stop()
	for {
		key, value, valid := next()
		if valid {
			fmt.Println(key, value)
		} else {
			break
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}
