package main

import (
	"fmt"
	"net/http"
	"time"
)

// go标准库web handler的写法
//
// 返回boy字符串
func getBoy(w http.ResponseWriter, r *http.Request) {
	defer func() { //如果用GIN，它会为每个Handler包一层recover。如果用标准http库，需要自己写recover
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi boy"))
}

// go标准库web handler的写法
//
// 返回girl字符串
func getGirl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi girl"))
}
