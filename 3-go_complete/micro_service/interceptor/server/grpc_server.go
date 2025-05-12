package main

import (
	"context"
	"dqq/micro_service/idl"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc "google.golang.org/grpc"
)

var limitAtServer = make(chan struct{}, 10) //瞬间并发度限制为10

// 计时+限流拦截器
func timerAndLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	begin := time.Now()
	limitAtServer <- struct{}{}
	resp, err := handler(ctx, req)
	fmt.Printf("concurrency %d\n", len(limitAtServer)) //打印瞬间并发度
	<-limitAtServer
	fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())
	return resp, err
}

func timerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	begin := time.Now()
	resp, err := handler(ctx, req)
	fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())
	return resp, err
}

func lLimitInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	limitAtServer <- struct{}{}
	resp, err := handler(ctx, req)
	fmt.Printf("concurrency %d\n", len(limitAtServer)) //打印瞬间并发度
	<-limitAtServer
	return resp, err
}

type MyServer struct {
	idl.UnimplementedHelloHttpServer
}

func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	return nil, nil
}

func (*MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: "hello"}
	return resp, nil
}

func main() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(
		// grpc.UnaryInterceptor(timerAndLimitInterceptor), //grpc.UnaryInterceptor只能使用一次，即server端只能会用一个拦截器
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			timerInterceptor,
			lLimitInterceptor,
		)), //让grpc_middleware帮你把多个拦截器封装成一个拦截器
	)
	// 注册服务
	idl.RegisterHelloServiceServer(server, new(MyServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run .\micro_service\interceptor\server\
