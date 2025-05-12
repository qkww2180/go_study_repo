package main

import (
	projectprepare "dqq/go/basic/project_prepare"
	"fmt"
	_ "net/http/pprof" //在线pprof

	_ "github.com/go-sql-driver/mysql" //注册mysql驱动
)

func InitLogger() {
	fmt.Println("init logger")
	fmt.Println("main是否匹配正则表达式", projectprepare.Reg.Match([]byte("hello123")))
}

func main() {
	projectprepare.CheckReg()
	InitLogger()
	InitDatabase()

	fmt.Println("server start")
}

func InitDatabase() {
	fmt.Println("init database")
}
