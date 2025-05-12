package test

import (
	"dqq/algorithm"
	"fmt"
	"math/rand"
	"testing"
)

var letterRunes = []rune("ABCDEFGHI")

// 生成随机字符串
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestLongestCommonSubstring(t *testing.T) {
	//人工观察一下结果是否正确
	s1 := "ABCDE"
	s2 := "BCECDE"
	cs1, result1 := algorithm.LongestCommonSubstring(s1, s2)
	cs2, result2 := algorithm.LongestCommonSubstringDP(s1, s2)
	cs3, result3 := algorithm.LongestCommonSubstringDP_WithSpaceON(s1, s2)
	cs4, result4 := algorithm.LongestCommonSubstringDP_WithSpaceON_StdCopy(s1, s2)
	if result1 != result2 || result2 != result3 || result3 != result4 {
		fmt.Println(s1)
		fmt.Println(s2)
		fmt.Println(cs1, cs2, cs3, cs4, result1, result2, result3, result4)
		t.Fail()
	}

	//随机选取1000个case，确保LongestCommonSubstring和LongestCommonSubstringDP的结果是相同的
	for i := 0; i < 1000; i++ {
		s1 = RandStringRunes(50) //随机生成长度为50的字符串
		s2 = RandStringRunes(50)
		cs1, result1 := algorithm.LongestCommonSubstring(s1, s2)
		cs2, result2 := algorithm.LongestCommonSubstringDP(s1, s2)
		cs3, result3 := algorithm.LongestCommonSubstringDP_WithSpaceON(s1, s2)
		cs4, result4 := algorithm.LongestCommonSubstringDP_WithSpaceON_StdCopy(s1, s2)
		if result1 != result2 || result2 != result3 || result3 != result4 {
			fmt.Println(s1)
			fmt.Println(s2)
			fmt.Println(cs1, cs2, cs3, cs4, result1, result2, result3, result4)
			t.Fail()
		}
	}
}

var (
	s1 = RandStringRunes(50)
	s2 = RandStringRunes(50)
)

func BenchmarkLongestCommonSubstring(b *testing.B) {
	for i := 0; i < b.N; i++ {
		algorithm.LongestCommonSubstring(s1, s2)
	}
}

func BenchmarkLongestCommonSubstringDP(b *testing.B) {
	for i := 0; i < b.N; i++ {
		algorithm.LongestCommonSubstringDP(s1, s2)
	}
}

func BenchmarkLongestCommonSubstringDP_WithSpaceON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		algorithm.LongestCommonSubstringDP_WithSpaceON(s1, s2)
	}
}
func BenchmarkLongestCommonSubstringDP_WithSpaceON_StdCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		algorithm.LongestCommonSubstringDP_WithSpaceON_StdCopy(s1, s2)
	}
}

func BenchmarkCopyRaw(b *testing.B) {
	arr1 := []byte(s1)
	arr2 := []byte(s2)
	n := len(arr1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < n; i++ {
			arr2[i] = arr1[i]
		}
	}
}

func BenchmarkCopyStd(b *testing.B) {
	arr1 := []byte(s1)
	arr2 := []byte(s2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(arr2, arr1)
	}
}

// go test ./algorithm/test -v -run=^TestLongestCommonSubstring -count=1
// go test ./algorithm/test -bench=^BenchmarkLongestCommonSubstring -run=^$ -count=1 -benchmem -benchtime=2s
// go test ./algorithm/test -bench=^BenchmarkCopy -run=^$ -count=1 -benchmem -benchtime=2s
