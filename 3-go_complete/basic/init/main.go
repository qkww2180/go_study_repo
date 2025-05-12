package main

import (
	"dqq/basic/init/db"
	"dqq/basic/init/logger"
	"fmt"
)

var (
	a int
	b int
)

func init() {
	fmt.Println("init")
	a = 9
	b = a + 9
}

var v func(int) (int, error)

var c = func() int {
	fmt.Println("set c")
	return 6
}()

func Init() {
	logger.InitLogger()
	db.InitDb()
}

func main() {
	fmt.Printf("a=%d, b=%d, c=%d\n", a, b, c)
}
