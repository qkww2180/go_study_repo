package main

import (
	"context"
	"dqq/micro_service/idl"
	"math"
	"net"

	"google.golang.org/grpc"
)

type EchoService struct {
	idl.UnimplementedEchoServer
}

func (*EchoService) Hello(ctx context.Context, request *idl.EchoRequest) (*idl.EchoResponse, error) {
	resp := idl.EchoResponse{Data: "hello"}
	return &resp, nil
}

func (*EchoService) StreamHello(server idl.Echo_StreamHelloServer) error {
	return server.SendMsg(&idl.EchoResponse{Data: "hello"})
}

func main2() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	//创建服务
	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(math.MaxInt),
		grpc.MaxSendMsgSize(math.MaxInt),
		grpc.InitialWindowSize(math.MaxInt32),
		grpc.InitialConnWindowSize(math.MaxInt32),
	)
	// 注册服务的具体实现
	idl.RegisterEchoServer(server, new(EchoService))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run ./micro_service/grpc/server
