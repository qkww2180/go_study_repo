package test

import (
	"dqq/algorithm"
	"dqq/algorithm/util"
	"fmt"
	"testing"
)

func TestLCS(t *testing.T) {
	s1 := []byte("AEMQ")
	s2 := []byte("BEQ")
	l1 := algorithm.LCS(s1, s2)
	l2 := algorithm.LCSTopDown(s1, s2)
	l3 := algorithm.LCSButtomUp(s1, s2)
	l4 := algorithm.LCSButtomUp_WithSpaceON(s1, s2)
	if l1 != l2 || l2 != l3 || l3 != l4 {
		t.Fail()
		fmt.Println(l1, l2, l3, l4)
	}
}

func BenchmarkLCS(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.LCS(s1, s2)
	}
}

func BenchmarkLCSTopDown(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.LCSTopDown(s1, s2)
	}
}

func BenchmarkLCSButtomUp(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.LCSButtomUp(s1, s2)
	}
}

func BenchmarkLCSButtomUp_WithSpaceON(b *testing.B) {
	s1 := []byte(util.RandStringRunes(10)) //生成长度为10的随机字符串
	s2 := []byte(util.RandStringRunes(10))
	b.ResetTimer() //生成s1,s2的耗时不包含在内
	for i := 0; i < b.N; i++ {
		algorithm.LCSButtomUp_WithSpaceON(s1, s2)
	}
}

// go test ./algorithm/test -v -run=^TestLCS$ -count=1
// go test ./algorithm/test -bench=^BenchmarkLCS -run=^$ -count=1 -benchmem -benchtime=2s

/**
BenchmarkLCS-8                              3720            793585 ns/op               0 B/op          0 allocs/op
BenchmarkLCSTopDown-8                    1642741              1403 ns/op            1344 B/op         12 allocs/op
BenchmarkLCSButtomUp-8                   2555702               926.9 ns/op          1344 B/op         12 allocs/op
BenchmarkLCSButtomUp_WithSpaceON-8       4913952               441.2 ns/op           192 B/op          2 allocs/op
**/
