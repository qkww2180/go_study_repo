package main

import (
	"context"
	"dqq/micro_service/idl"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpc_port = ":5678"
	http_port = ":5679"
)

type GreetService struct {
	idl.UnimplementedHelloHttpServer
}

func (s *GreetService) Greeting(ctx context.Context, request *idl.GreetRequest) (*idl.GreetResopnse, error) {
	resp := &idl.GreetResopnse{
		Message: "Hello " + request.Name,
	}
	return resp, nil
}

// 在2个端口上分别启动了grpc和http服务
func main1() {
	listener, err := net.Listen("tcp", grpc_port)
	if err != nil {
		panic(err)
	}
	grpcHandler := grpc.NewServer() //grpcHandler实现了http.Handler接口
	//绑定服务的实现
	idl.RegisterHelloHttpServer(grpcHandler, new(GreetService))
	go func() {
		//放到子协程里启动grpc服务
		grpcHandler.Serve(listener)
	}()

	// 连接grpc服务
	conn, err := grpc.Dial(grpc_port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	gwHandler := runtime.NewServeMux() //实现了http.Handler接口
	//把post到/golang/hello上的请求转到grpc connection上
	if err := idl.RegisterHelloHttpHandler(context.Background(), gwHandler, conn); err != nil {
		panic(err)
	}

	// 启动http服务
	if err := http.ListenAndServe(http_port, gwHandler); err != nil {
		fmt.Println(err)
	}
}

// go run .\micro_service\gateway\server\
