package test

import (
	"dqq/algorithm"
	"testing"
)

var binaryTree *algorithm.BNode

func init() {
	binaryTree = &algorithm.BNode{Value: 5}
	n15 := &algorithm.BNode{Value: 15}
	n10 := &algorithm.BNode{Value: 10}
	n20 := &algorithm.BNode{Value: 20}
	n30 := &algorithm.BNode{Value: 30}
	n62 := &algorithm.BNode{Value: 62}
	n49 := &algorithm.BNode{Value: 49}
	binaryTree.LeftChild = n15
	binaryTree.RightChild = n10
	n15.LeftChild = n20
	n15.RightChild = n30
	n10.LeftChild = n62
	n10.RightChild = n49
}

func TestPreOrder(t *testing.T) {
	binaryTree.PreOrder()
}

func TestPostOrder(t *testing.T) {
	binaryTree.PostOrder()
}
func TestMiddleOrder(t *testing.T) {
	binaryTree.MiddleOrder()
}

// go test ./algorithm/test -v -run=^TestPreOrder$ -count=1
// go test ./algorithm/test -v -run=^TestPostOrder$ -count=1
// go test ./algorithm/test -v -run=^TestMiddleOrder$ -count=1
