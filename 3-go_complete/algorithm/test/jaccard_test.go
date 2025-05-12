package test

import (
	"dqq/algorithm"
	"fmt"
	"testing"

	"golang.org/x/exp/slices" //golang.org/x下的代码处于试验阶段，成熟后会放到标准库下
)

func TestJaccardTimeConsuming(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器"}
	fmt.Println(algorithm.JaccardTimeConsuming(l1, l2))
}

func TestJaccardForSorted(t *testing.T) {
	l1 := []string{"go", "分布式", "mysql", "搞笑", "并发编程", "服务器"}
	l2 := []string{"c#", "AI", "mysql", "篮球", "并发编程", "服务器"}
	slices.Sort(l1)
	slices.Sort(l2)
	fmt.Println(algorithm.JaccardForSorted(l1, l2))
}

// go test ./algorithm/test -v -run=^TestJaccardTimeConsuming$ -count=1
// go test ./algorithm/test -v -run=^TestJaccardForSorted$ -count=1
