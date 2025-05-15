package main

import (
	"dqq/micro_service/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func greet(userId int) string {
	url := "http://127.0.0.1:5678/greet"
	begin := time.Now()
	traceId := util.RandStringRunes(10) //生成随机字符串
	defer func() {
		log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, strconv.Itoa(userId), url, begin.UnixNano(), time.Since(begin).Nanoseconds())
	}()

	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("new h_http request failed: %s", err)
		return ""
	}
	request.Header.Set("trace_id", traceId)
	request.Header.Set("user_id", strconv.Itoa(userId))

	response, err := client.Do(request)
	if err != nil {
		log.Printf("call h_http server failed: %s", err)
		return ""
	}
	defer response.Body.Close()
	bs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("read h_http response failed: %s", err)
		return ""
	}
	return string(bs)
}

func main2() {
	userId := 8
	fmt.Println(greet(userId))
}

// go run .\micro_service\trace\self_trace\rpc_tracer\client\
