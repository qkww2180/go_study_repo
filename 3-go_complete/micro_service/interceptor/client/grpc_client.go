package main

import (
	"context"
	"dqq/micro_service/idl"
	"fmt"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var limitAtClient = make(chan struct{}, 10) //瞬间并发度限制为10

// 计时拦截器
func timerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())
	return err
}

// 限流拦截器
func limitInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	limitAtClient <- struct{}{}
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("concurrency %d\n", len(limitAtClient)) //打印瞬间并发度
	<-limitAtClient
	return err
}

func main() {
	// 连接到GRPC服务端
	conn, err := grpc.Dial("127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()), //无需使用安全传输
		//在GRPC客户端使用拦截器
		grpc.WithChainUnaryInterceptor(timerInterceptor, limitInterceptor), //不定长参数，第一个拦截器在最外层，最后一个拦截器最靠近真实的业务调用
	)
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}
	defer conn.Close()

	client := idl.NewHelloServiceClient(conn)

	// 执行RPC调用并打印收到的响应数据，指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, &idl.HelloRequest{})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return
	}
	fmt.Println(resp.Greeting)
}

// go run .\micro_service\interceptor\client\
