package main

import (
	student_service "dqq/micro_service/grpc"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
)

func getStudent(w http.ResponseWriter, r *http.Request) {
	defer func() { //GIN框架已经为每个handler写了recover()，如果使用http标准库还是需要自己写recover()
		panicInfo := recover() //panicInfo是any类型，即传给panic()的参数
		if panicInfo != nil {
			fmt.Println(panicInfo)
		}
	}()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("read request failed"))
		return
	}
	r.Body.Close()
	var request student_service.StudentID
	err = sonic.Unmarshal(bs, &request)
	if err != nil {
		w.Write([]byte("request is not valid json"))
		return
	}
	student := &student_service.Student{Name: "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
		Id:        request.Id,
	}
	bs, err = sonic.Marshal(student)
	if err != nil {
		w.Write([]byte("marshal response failed"))
		return
	}
	w.Write(bs)
}

func main() {
	http.Handle("/", http.HandlerFunc(getStudent))
	if err := http.ListenAndServe("127.0.0.1:5679", nil); err != nil {
		fmt.Println(err)
	}
}

// go run .\micro_service\grpc\web_server\
