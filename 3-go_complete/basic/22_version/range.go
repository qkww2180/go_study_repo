package main

import (
	"fmt"
	"sync"
)

func main1() {
	wg := sync.WaitGroup{}
	values := []string{"a", "b", "c"}
	wg.Add(len(values))
	for _, v := range values { //每次for range 都会为v重新申请一块内存
		fmt.Printf("%p\n", &v)
		go func() {
			defer wg.Done()
			fmt.Println(v)
		}()
	}
	wg.Wait()

	for i := range 3 { //v1.22 新的遍历方式
		fmt.Println(i)
	}
}
