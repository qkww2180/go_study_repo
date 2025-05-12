package test

import (
	"container/list"
	"dqq/algorithm"
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	lst := new(algorithm.DoubleList[int])
	lst.PushBack(1)  // 1
	lst.PushBack(2)  // 1 -> 2
	lst.PushFront(3) // 3 -> 1 -> 2
	lst.PushFront(4) // 4 -> 3 -> 1 -> 2

	third := lst.Get(2)        //第3个元素是1
	lst.InsertAfter(8, third)  // 4 -> 3 -> 1 -> 8 -> 2
	lst.InsertBefore(9, third) // 4 -> 3 -> 9 -> 1 -> 8 -> 2

	fmt.Printf("链表中共%d个元素\n", lst.Length)
	lst.Traverse()
	lst.ReverseTraverse()
}

func TestStdList(t *testing.T) {
	lst := list.New()
	lst.PushBack(1)  // 1
	lst.PushBack(2)  // 1 -> 2
	lst.PushFront(3) // 3 -> 1 -> 2
	lst.PushFront(4) // 4 -> 3 -> 1 -> 2

	third := lst.Front().Next().Next() //第3个元素是1
	lst.InsertAfter(8, third)          // 4 -> 3 -> 1 -> 8 -> 2
	lst.InsertBefore(9, third)         // 4 -> 3 -> 9 -> 1 -> 8 -> 2
	lst.Remove(third)                  // 4 -> 3 -> 9 -> 8 -> 2

	fmt.Printf("链表中共%d个元素\n", lst.Len())
	algorithm.TraversList(lst)
	algorithm.ReverseList(lst)
}

// go test ./algorithm/test -v -run=^TestList$ -count=1
// go test ./algorithm/test -v -run=^TestStdList$ -count=1
