package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/trace"
	"dqq/micro_service/trace/open_trace/rpc_tracer/common"
	"dqq/micro_service/util"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type MyServer struct {
	idl.UnimplementedHelloHttpServer
}

func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	return &idl.LoginResponse{Token: util.RandStringRunes(10)}, nil
}

func (*MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: "hello"}
	return resp, nil
}

func InitGrpcJaeger() {
	var jaeger opentracing.Tracer
	var err error
	// 设置全局tracer
	jaeger, closer, err = trace.NewJaegerTracer("my_grpc_server", "127.0.0.1:6831") //需要先启动jaeger
	if err != nil {
		panic(err)
	}
	// defer closer.Close()  需要在接收到kill信号时调用Close()
	opentracing.SetGlobalTracer(jaeger)
}

func main2() {
	InitGrpcJaeger()
	defer closer.Close()
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(common.ServerTraceInterceptor))
	// 绑定服务实现
	idl.RegisterHelloServiceServer(server, new(MyServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run .\micro_service\trace\open_trace\rpc_tracer\server\
