package algorithm

import (
	"cmp"
	"fmt"
)

type ListNode[T cmp.Ordered] struct {
	Value T
	Prev  *ListNode[T]
	Next  *ListNode[T]
}

// 双向链表
type DoubleList[T cmp.Ordered] struct {
	Head   *ListNode[T]
	Tail   *ListNode[T]
	Length int
}

// 正序遍历整个List
func (list *DoubleList[T]) Traverse() {
	curr := list.Head //从首部开始
	for curr != nil {
		fmt.Printf("%v ", curr.Value)
		curr = curr.Next //沿着Next往后走
	}
	fmt.Println()
}

// 逆序遍历整个List
func (list *DoubleList[T]) ReverseTraverse() {
	curr := list.Tail
	for curr != nil {
		fmt.Printf("%v ", curr.Value)
		curr = curr.Prev //沿着Prev往前走
	}
	fmt.Println()
}

// 向尾部追加一个元素。O(1)
func (list *DoubleList[T]) PushBack(x T) {
	node := &ListNode[T]{Value: x} //把x封装成Node，//node的Prev、Next有待赋值，现在还是nil
	tail := list.Tail
	if tail == nil { //list为空
		list.Head = node //node既是首，也是尾
		list.Tail = node
	} else {
		tail.Next = node //tail是原先的尾节点
		node.Prev = tail
		list.Tail = node //List有了新的tail
	}
	list.Length += 1
}

// 向首部追加一个元素。O(1)
func (list *DoubleList[T]) PushFront(x T) {
	node := &ListNode[T]{Value: x} //把x封装成Node，//node的Prev、Next有待赋值，现在还是nil
	head := list.Head
	if head == nil { //list为空
		list.Head = node //node既是首，也是尾
		list.Tail = node
	} else {
		head.Prev = node //head是原先的首部
		node.Next = head
		list.Head = node //List有了新的head
	}
	list.Length += 1
}

// 获取第idx个元素。O(N)
func (list *DoubleList[T]) Get(idx int) *ListNode[T] {
	if list.Length <= idx {
		return nil
	}
	curr := list.Head
	for i := 0; i < idx; i++ {
		curr = curr.Next
	}
	return curr
}

// 在n1后面添加一个元素x。O(N)
func (list *DoubleList[T]) InsertAfter(x T, n1 *ListNode[T]) {
	n2 := &ListNode[T]{Value: x} //node的Prev、Next有待赋值
	if n1.Next != nil {          //prevNode本来是尾元素
		n3 := n1.Next //获取prevNode的下一个元素，node就要插入到它俩之间
		n3.Prev = n2  //插入一个node，原有的Prev、Next哪些会受到影响，要考虑清楚
		n2.Next = n3
	} else {
		list.Tail = n2
	}
	n1.Next = n2
	n2.Prev = n1
	list.Length += 1
}

// 在n3前面添加一个元素x。O(N)
func (list *DoubleList[T]) InsertBefore(x T, n3 *ListNode[T]) {
	n2 := &ListNode[T]{Value: x} //node的Prev、Next有待赋值
	if n3.Prev != nil {          //nextNode本来是首元素
		n1 := n3.Prev
		n1.Next = n2
		n2.Prev = n1
	} else {
		list.Head = n2
	}
	n3.Prev = n2
	n2.Next = n3
	list.Length += 1
}
