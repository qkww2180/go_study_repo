package main

import (
	"fmt"
	"net/http"
	"time"
)

func welcome(ch chan<- string) {
	time.Sleep(2 * time.Second) //故意慢一点，2秒后才返回结果
	ch <- "welcome"
}

func handler(w http.ResponseWriter, req *http.Request) {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	ctx := req.Context() //取得request的context
	blocker := make(chan string, 1)
	go welcome(blocker)
	select {
	case v := <-blocker:
		fmt.Fprint(w, v)
	case <-ctx.Done(): //超时后client会撤销请求，触发ctx.cancel()，从而关闭Done()管道
		err := ctx.Err()            //如果发生Done（管道被关闭），Err返回Done的原因，可能是被Cancel了，也可能是超时了
		fmt.Println("server:", err) //context canceled
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:5678", nil)
}

// go run ./basic/timeout/server
