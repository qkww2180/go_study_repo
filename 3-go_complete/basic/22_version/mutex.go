package main

import (
	"fmt"
	"net/http"
)

func main3() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) { //v1.22可以直接在路由中指定允许的请求方法
		// if r.Method != "GET" {   //v1.22之前需要自己限制请求方法
		// 	fmt.Fprint(w, "warn: 只支持GET方法")
		// } else {
		fmt.Fprint(w, "你好 "+r.PathValue("name")) //从restful风格的url中获取参数
		// }
	})

	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
}
