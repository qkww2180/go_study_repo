package main

import (
	"fmt"
	"net/http"
)

// 使用go标准库实现跟GIN中间件一样的功能

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler //mux通常表示路由策略
}

func NewRouter() *Router {
	return &Router{
		middlewareChain: make([]middleware, 0, 10),
		mux:             make(map[string]http.Handler, 10),
	}
}

func (router *Router) Use(m middleware) {
	router.middlewareChain = append(router.middlewareChain, m)
}

func (router *Router) Add(path string, handler http.Handler) {
	var mergedHandler = handler
	for i := len(router.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = router.middlewareChain[i](mergedHandler) //中间件层层嵌套
	}
	router.mux[path] = mergedHandler
}

// Router实现了ServeHTTP方法，所以说Router实现了http.Handler接口
func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.RequestURI()
	//在Handler内部实现路由
	if handler, exists := router.mux[requestPath]; !exists {
		http.NotFoundHandler().ServeHTTP(w, r) //调用标准库自带的NotFoundHandler
	} else {
		handler.ServeHTTP(w, r)
	}
}

func main2() {
	router := NewRouter()
	router.Use(limitMiddleWare)
	router.Use(timeMiddleWare)
	//以下演示了2个路径（还可以更多），每个路径都使用相同的middlewareChain
	router.Add("/", http.HandlerFunc(getBoy))
	router.Add("/home", http.HandlerFunc(getGirl))

	// 用http.Handle()实现路由，这样http.ListenAndServe()的第二个参数Handler可以传nil
	// for path, handler := range router.mux {
	// 	http.Handle(path, handler)
	// }

	if err := http.ListenAndServe("127.0.0.1:5678", router); err != nil {
		fmt.Println(err)
	}
}
