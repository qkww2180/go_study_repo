package main

import "net/http"

func main24() {
	var f1, f2 func(a, b int, c string, d bool) (int, bool)

	f1 = func(a, b int, c string, d bool) (int, bool) {
		return a + b, d && c == "abc"
	}
	// 没给f2赋值，此时f2为nil，如果调用f2会发生空指针异常

	i, p := f1(1, 2, "3", false)
	_, _ = i, p

	type ConnectionPool struct {
		Servers      []string
		LoadBalancer func(a, b int, c string, d bool) (int, bool) //成员变量是接口类型
	}

	cp := ConnectionPool{
		Servers:      []string{"127.0.0.1:1234", "127.0.0.1:5678"},
		LoadBalancer: f2,
	}
	_ = cp
	// cp.LoadBalancer(1, 2, "3", false)  // nil pointer dereference

	type ConnectionPool2 struct {
		Servers      []string
		LoadBalancer FT //成员变量是接口类型
	}
	cp2 := ConnectionPool2{
		Servers:      []string{"127.0.0.1:1234", "127.0.0.1:5678"},
		LoadBalancer: FT(f2),
	}
	_ = cp2

	transport("BJ", "SH", FT(f1)) //函数f1并没有实现Car接口，但FT类型实现了

	h := func(w http.ResponseWriter, r *http.Request) {}
	http.ListenAndServe(":8080", http.HandlerFunc(h)) //函数h并没有实现http.Handler接口，但http.HandlerFunc类型实现了
}

type FT func(a, b int, c string, d bool) (int, bool)

func (FT) move(src string, dest string) (int, error) {
	return 0, nil
}

func (FT) whistle(n int) int {
	return 0
}
