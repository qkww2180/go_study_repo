package test

import (
	"dqq/algorithm"
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	tree := new(algorithm.TrieTree)
	tree.AddTerm("分散")
	tree.AddTerm("分散精力")
	tree.AddTerm("分散投资")
	tree.AddTerm("分布式")
	tree.AddTerm("工程")
	tree.AddTerm("工程师")

	terms := tree.Retrieve("分散")
	fmt.Println(terms)
	terms = tree.Retrieve("人工")
	fmt.Println(terms)
}

// go test ./algorithm/test -v -run=^TestTrie$ -count=1
