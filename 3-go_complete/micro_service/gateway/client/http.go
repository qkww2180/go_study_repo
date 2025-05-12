package main

import (
	"context"
	"dqq/micro_service/idl"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpc_port = ":5678"
	http_port = ":5679"
)

var request = &idl.GreetRequest{
	Name: "大乔乔",
}

func grpcClient() {
	// 连接到GRPC服务端
	conn, err := grpc.Dial(grpc_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}
	defer conn.Close()
	client := idl.NewHelloHttpClient(conn)

	// 执行RPC调用并打印收到的响应数据，指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Greeting(ctx, request)
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return
	} else {
		fmt.Println(resp.Message)
	}
}

func httpClient() {
	client := &http.Client{}

	reqStr, err := sonic.MarshalString(request)
	if err != nil {
		return
	}
	reader := strings.NewReader(reqStr)

	// 在.proto文件里指定的请求方法是post，路径是"/golang/hello"
	if resp, err := client.Post("http://127.0.0.1"+http_port+"/golang/hello", "application/json", reader); err != nil { //请求是json格式
		panic(err)
	} else {
		defer client.CloseIdleConnections()
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		bs, _ := io.ReadAll(resp.Body)
		var response idl.GreetResopnse
		sonic.Unmarshal(bs, &response)
		fmt.Println(response.Message)
	}
}

func main1() {
	grpcClient()
	httpClient()
}

// go run .\micro_service\gateway\client\
