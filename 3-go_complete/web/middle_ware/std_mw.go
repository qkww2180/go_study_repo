package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var limitCh = make(chan struct{}, 100) //最多并发处理100个请求

func timeMiddleWare(next http.Handler) http.Handler {
	//通过HandlerFunc把一个func(rw http.ResponseWriter, r *http.Request)函数转为Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(rw, r)
		timeElapsed := time.Since(begin)
		log.Printf("request %s use %d ms\n", r.URL.Path, timeElapsed.Milliseconds())
	})
}

func limitMiddleWare(next http.Handler) http.Handler {
	//通过HandlerFunc返回一个Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		limitCh <- struct{}{} //并发度达到100时就会阻塞
		log.Printf("concurrence %d\n", len(limitCh))
		next.ServeHTTP(rw, r)
		<-limitCh
	})
}

func main1() {
	//用http.Handle实现路由
	http.Handle("/", timeMiddleWare(limitMiddleWare(http.HandlerFunc(getBoy))))      //中间件层层嵌套
	http.Handle("/home", timeMiddleWare(limitMiddleWare(http.HandlerFunc(getGirl)))) //跟上面一行存在重复代码，不够优雅

	//http.ListenAndServe()的第二个参数handler可以传nil
	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		fmt.Println(err)
	}
}
