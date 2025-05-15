package g_database

import (
	"log"
	rand "math/rand/v2"
	"os"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// 生成随机字符串
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}
