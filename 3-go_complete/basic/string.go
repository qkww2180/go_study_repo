package main

import (
	"fmt"
	"unicode/utf8"
)

func main9() {
	s := "go语言"
	arr := []byte(s)
	fmt.Println(arr)
	fmt.Println(arr[2], s[2])
	arr[2] = 9
	//s[2] = 9 //字符串不能修改

	fmt.Println(len(s), len(arr))

	brr := []rune(s) //一个汉字是一个rune
	fmt.Println(brr)
	fmt.Printf("%d %d %s %s\n", brr[2], s[2], string(brr[2]), string(s[2])) //数字可以直接转string
	fmt.Println(utf8.RuneCountInString(s))                                  //查看string里有几个rune
}
