package util

import "math/rand"

// 是否为大写英文字母
func IsASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// 大小写英文字母互相转换
func UpperLowerExchange(c byte) byte {
	return c ^ ' '
}

// 驼峰转为蛇形
func Camel2Snake(s string) string {
	if len(s) == 0 {
		return ""
	}
	t := make([]byte, 0, len(s)+4)
	if IsASCIIUpper(s[0]) { //首字母为大写时，转为小写
		t = append(t, UpperLowerExchange(s[0]))
	} else {
		t = append(t, s[0])
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if IsASCIIUpper(c) { //碰到大写，则转为小写，并在前面加个_
			t = append(t, '_', UpperLowerExchange(c))
		} else {
			t = append(t, c)
		}
	}
	return string(t)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 生成随机字符串
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
