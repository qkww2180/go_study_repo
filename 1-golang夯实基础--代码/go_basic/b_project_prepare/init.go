package projectprepare

import (
	"fmt"
	"regexp"
)

var Reg *regexp.Regexp

// 在init()函数内不要依赖外部的其他任何变量
func init() {
	fmt.Printf("a=%d\n", a)
	var err error
	Reg, err = regexp.Compile(`\d+`)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("init Reg success")
	}
}

func CheckReg() {
	fmt.Println("init是否匹配正则表达式", Reg.Match([]byte("hello")))
	_ = b
	_ = c
}

func gn() int {
	fmt.Println("gn called")
	return 999
}

var a int = 9
var b int = gn()
var c int = func() int {
	fmt.Println("anonymous called")
	return 999
}()
