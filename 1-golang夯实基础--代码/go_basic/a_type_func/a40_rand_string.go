package main

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

var letterCollection = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY/_+大乔中国")

// 生成长度为n的随机字符串
func RandString(n int) string {
	// rect := make([]rune, 0, n)
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		index := rand.IntN(len(letterCollection))
		// rect = append(rect, letterCollection[index])
		sb.WriteRune(letterCollection[index])
	}
	// return string(rect)
	return sb.String()
}

func main40() {
	fmt.Println(RandString(10))
	fmt.Println(RandString(10))
	fmt.Println(RandString(10))
	fmt.Println(RandString(10))
}
