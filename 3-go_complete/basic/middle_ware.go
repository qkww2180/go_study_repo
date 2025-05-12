package main

import (
	"fmt"
	"net/http"
)

type Handler interface {
	Serve(int) int
}

type run int
type HandlerFunc func(int, int) int

func (r run) Serve(arg int) int {
	return arg * 2
}

func (h HandlerFunc) Serve(arg int) int {
	return h(arg, 4) * 2
}

// func init() {
// 	var r run = 8
// 	fmt.Printf("%d\n", r)

// 	f := func(a int, b int) int {
// 		return a + b
// 	}

// 	var a HandlerFunc = HandlerFunc(f)
// 	fmt.Printf("%d\n", a(2, 5))

// 	var b Handler = a
// 	fmt.Printf("%d\n", b.Serve(3))
// }

func getSex() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Male\n"))
		})
}

func getName() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Tom\n"))
		})
}

func withMiddleWare1(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("这是中间件1的结果\n"))
			next.ServeHTTP(w, r)
		})
}

func withMiddleWare2(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("这是中间件2的结果\n"))
			next.ServeHTTP(w, r)
		})
}

func main18() {
	http.Handle("/a", getSex())
	http.Handle("/b", getName())
	http.Handle("/c", withMiddleWare1(getSex()))
	http.Handle("/d", withMiddleWare1(getName()))
	http.Handle("/e", withMiddleWare1(withMiddleWare2(withMiddleWare1(getName()))))

	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		fmt.Println(err)
	}
}

// go run ./basic
