package main

import (
	"html/template"
	"net/http"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {
	tmpl := template.Must(template.ParseFiles("./web/template/home.html")) //解析模板文件
	//定义路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "这里是标题",
			Todos: []Todo{
				{Title: "任务一，未完成", Done: false},
				{Title: "任务二，完成", Done: true},
				{Title: "任务三，完成", Done: true},
			},
		}
		tmpl.Execute(w, data) //向模板里填充具体数据
	})
	//启动http server
	http.ListenAndServe("localhost:5678", nil) //localhost不需要走网卡，127.0.0.1需要走网卡
}

// go run ./web/template
