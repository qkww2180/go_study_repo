package algorithm

import (
	"container/list"
)

// LRU缓存。不支持并发
type LRUCache struct {
	cache map[int]string //存储业务数据
	lst   list.List      //辅助的链表，用于决定淘汰谁
	cap   int            //缓存容量的上限
}

func NewLRUCache(cap int) *LRUCache {
	lru := new(LRUCache)
	lru.cap = cap
	lru.cache = make(map[int]string, cap)
	lru.lst = list.List{}
	return lru
}

// 向缓存中添加元素(添加之前先通过调Get确保key不在缓存里)。时间复杂度O(1)
func (lru *LRUCache) Add(key int, value string) {
	if len(lru.cache) == lru.cap { //刚刚到达缓存容量上限
		//先从缓存中淘汰一个元素
		back := lru.lst.Back()
		delete(lru.cache, back.Value.(int)) //interface {} is nil, not int
		lru.lst.Remove(back)
	}
	//把key value放到缓存中去
	lru.cache[key] = value
	lru.lst.PushFront(key)
}

// 从辅助链表中找一个元素。时间复杂度O(N)
func (lru *LRUCache) find(key int) *list.Element {
	if lru.lst.Len() == 0 {
		return nil
	}
	head := lru.lst.Front()
	for {
		if head == nil {
			break
		}
		if head.Value.(int) == key {
			return head
		} else {
			head = head.Next()
		}
	}
	return nil
}

// 从缓存中查找元素。时间复杂度O(N)
func (lru *LRUCache) Get(key int) (string, bool) {
	value, exists := lru.cache[key]
	ele := lru.find(key)
	if ele != nil {
		lru.lst.MoveToFront(ele)
	}
	return value, exists
}
