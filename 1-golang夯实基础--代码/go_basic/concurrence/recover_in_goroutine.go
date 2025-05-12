package concurrence

import (
	"strconv"
	"time"
)

var table = []int{1, 2, 3, 4, 5, 6, 7}

func GetHandler(index int) int {
	defer func() {
		recover() //recover()只能捕获本协程内的panic
	}()
	return table[index]
}

func SetHandler(index int, div string) {
	defer func() {
		recover() //recover()只能捕获本协程内的panic
	}()
	i, _ := strconv.Atoi(div)
	table[index] = 10 / i
}

func ServiceMain() {
	go GetHandler(10)
	go SetHandler(0, "7s")
	time.Sleep(time.Second)
}
