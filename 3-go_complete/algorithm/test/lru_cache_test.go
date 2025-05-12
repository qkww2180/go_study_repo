package test

import (
	"dqq/algorithm"
	"fmt"
	"strconv"
	"testing"
)

func TestLRUCache(t *testing.T) {
	lru := algorithm.NewLRUCache(10) //缓存容量为10
	for i := 0; i < 10; i++ {        //填满缓存
		lru.Add(i, strconv.Itoa(i)) //9 8 7 6 5 4 3 2 1 0
	}

	for i := 0; i < 10; i += 2 { //访问偶数元素。被访问的元素会放到链表的首部
		lru.Get(i) //8 6 4 2 0 9 7 5 3 1
	}

	for i := 10; i < 15; i++ { //再添加5个新元素。新添加的元素会放到链表的首部
		lru.Add(i, strconv.Itoa(i)) //14 13 12 11 10 8 6 4 2 0
	}

	for i := 0; i < 10; i++ { //检查缓存中还有没有最初的那10个元素
		_, exists := lru.Get(i)
		fmt.Printf("key %d exists %t\n", i, exists) //9 7 5 3 1不存在，8 6 4 2 0存在
	}
}

// go test ./algorithm/test -v -run=^TestLRUCache$ -count=1