package main

import (
	"flag"
	"fmt"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/chat_room/home.html") //请求根目录时直接返回一个html页面
}

func main() {
	port := flag.String("port", "5657", "h_http service port") //如果命令行不指定port参数，则默认为5657
	flag.Parse()                                               //解析命令行输入的port参数
	hub := NewHub()
	go hub.Run()
	//定义路由，注册每种请求对应的处理函数
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/chat", func(rw http.ResponseWriter, r *http.Request) {
		ServeWs(hub, rw, r)
	})
	fmt.Printf("h_http serve on port %s\n", *port)
	if err := http.ListenAndServe("127.0.0.1:"+*port, nil); err != nil { //如果启动成功，该行会一直阻塞，hub.run()会一直运行。由于上面已经定义了路由，所以这里handler参数可以传nil
		fmt.Printf("start h_http service error: %s\n", err)
	}
}

// go run ./web/chat_room --port 5678
