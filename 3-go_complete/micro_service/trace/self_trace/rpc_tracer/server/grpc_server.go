package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"log"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func traceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	begin := time.Now()
	var traceId string
	var userId string

	//从Client端传过来的ctx里抽取出metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if value, exists := md["trace_id"]; exists {
			traceId = value[0]
		}
		if value, exists := md["user_id"]; exists {
			userId = value[0]
		}
	}

	resp, err := handler(ctx, req)

	log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, userId, info.FullMethod, begin.UnixNano(), time.Since(begin).Nanoseconds())
	return resp, err
}

type MyServer struct {
	idl.UnimplementedHelloHttpServer
}

func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	return &idl.LoginResponse{Token: util.RandStringRunes(10)}, nil
}

func (*MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	begin := time.Now()
	var traceId string
	var userId string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if value, exists := md["trace_id"]; exists {
			traceId = value[0]
		}
		if value, exists := md["user_id"]; exists {
			userId = value[0]
		}
	}
	defer func() {
		log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, userId, "SayHello", begin.UnixNano(), time.Since(begin).Nanoseconds())
	}()

	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: "hello"}
	return resp, nil
}

func main1() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(traceInterceptor))
	// 绑定服务实现
	idl.RegisterHelloServiceServer(server, new(MyServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run .\micro_service\trace\self_trace\rpc_tracer\server\
