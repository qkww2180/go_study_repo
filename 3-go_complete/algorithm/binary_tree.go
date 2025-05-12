package algorithm

import "fmt"

type BNode struct {
	Value      any
	LeftChild  *BNode
	RightChild *BNode
}

func (node *BNode) PreOrder() {
	if node == nil {
		return
	}
	fmt.Printf("%v ", node.Value)
	node.LeftChild.PreOrder()
	node.RightChild.PreOrder()
}

func (node *BNode) PostOrder() {
	if node == nil {
		return
	}
	node.LeftChild.PostOrder()
	node.RightChild.PostOrder()
	fmt.Printf("%v ", node.Value)
}

func (node *BNode) MiddleOrder() {
	if node == nil {
		return
	}
	node.LeftChild.MiddleOrder()
	fmt.Printf("%v ", node.Value)
	node.RightChild.MiddleOrder()
}
