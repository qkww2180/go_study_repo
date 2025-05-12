package algorithm

import (
	"container/list"
	"fmt"
)

// 正序遍历整个List
func TraversList(lst *list.List) {
	head := lst.Front() //取到首元素
	for head.Next() != nil {
		fmt.Printf("%v ", head.Value)
		head = head.Next()
	}
	fmt.Println(head.Value)
}

// 逆序遍历整个List
func ReverseList(lst *list.List) {
	tail := lst.Back() //取到末元素
	for tail.Prev() != nil {
		fmt.Printf("%v ", tail.Value)
		tail = tail.Prev()
	}
	fmt.Println(tail.Value)
}
