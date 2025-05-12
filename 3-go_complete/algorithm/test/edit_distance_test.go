package test

import (
	"dqq/algorithm"
	"dqq/algorithm/util"
	"fmt"
	"testing"
)

func TestEditDistance(t *testing.T) {
	s1 := []byte("word")
	s2 := []byte("world")
	if algorithm.EditDistance(s1, s2) != 1 {
		fmt.Printf("EditDistance fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceTopDown(s1, s2) != 1 {
		fmt.Printf("EditDistanceTopDown fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceButtomUp(s1, s2) != 1 {
		fmt.Printf("EditDistanceButtomUp fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceButtomUp_WithSpaceON(s1, s2) != 1 {
		fmt.Printf("EditDistanceButtomUp_WithSpaceON fail for %s %s\n", s1, s2)
		t.Fail()
	}

	s1 = []byte("horse")
	s2 = []byte("ros")
	if algorithm.EditDistance(s1, s2) != 3 {
		fmt.Printf("EditDistance fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceTopDown(s1, s2) != 3 {
		fmt.Printf("EditDistanceTopDown fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceButtomUp(s1, s2) != 3 {
		fmt.Printf("EditDistanceButtomUp fail for %s %s\n", s1, s2)
		t.Fail()
	}
	if algorithm.EditDistanceButtomUp_WithSpaceON(s1, s2) != 3 {
		fmt.Printf("EditDistanceButtomUp_WithSpaceON fail for %s %s\n", s1, s2)
		t.Fail()
	}
}

func BenchmarkEditDistance(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.EditDistance(s1, s2)
	}
}

func BenchmarkEditDistanceTopDown(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.EditDistanceTopDown(s1, s2)
	}
}

func BenchmarkEditDistanceButtomUp(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.EditDistanceButtomUp(s1, s2)
	}
}

func BenchmarkEditDistanceButtomUp_WithSpaceON(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.EditDistanceButtomUp_WithSpaceON(s1, s2)
	}
}

// go test ./algorithm/test -v -run=^TestEditDistance$ -count=1
// go test ./algorithm/test -bench=^BenchmarkEditDistance -run=^$ -count=1 -benchmem -benchtime=2s

/**
goos: windows
goarch: amd64
pkg: upgrading/algorithm/test
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkEditDistance-8                              100          24605215 ns/op               0 B/op              0 allocs/op

BenchmarkEditDistanceTopDown-8                   1462236              1659 ns/op            1344 B/op             12 allocs/op

BenchmarkEditDistanceButtomUp-8                  2471469              1060 ns/op            1344 B/op         12 allocs/op
BenchmarkEditDistanceButtomUp_WithSpaceON-8      4959926               490.2 ns/op           192 B/op          2 allocs/op
**/
