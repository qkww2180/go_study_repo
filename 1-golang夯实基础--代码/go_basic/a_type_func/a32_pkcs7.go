package main

import (
	"bytes"
	"errors"
	"fmt"
)

func Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - srcLen%blockSize
	// rect := src
	// for i := 0; i < padLen; i++ {
	// 	rect = append(rect, byte(padLen))
	// }
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	rect := append(src, pad...)
	return rect
}

func Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	if srcLen%blockSize != 0 || srcLen <= blockSize {
		return nil, errors.New("参数有问题")
	}
	padLen := int(src[srcLen-1])
	return src[0 : srcLen-padLen], nil
}

func main32() {
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(Padding(src, 8))
	fmt.Println(Unpadding(Padding(src, 8), 8))
	src = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Println(Padding(src, 8))
	fmt.Println(Unpadding(Padding(src, 8), 8))
}
