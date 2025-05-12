package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/trace"
	"dqq/micro_service/trace/open_trace/rpc_tracer/common"
	"fmt"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func InitGrpcJaeger() {
	var jaeger opentracing.Tracer
	var err error
	// 设置全局tracer
	jaeger, closer, err = trace.NewJaegerTracer("my_grpc_client", "127.0.0.1:6831") //需要先启动jaeger
	if err != nil {
		panic(err)
	}
	// defer closer.Close()  需要在接收到kill信号时调用Close()
	opentracing.SetGlobalTracer(jaeger)
}

// RPC调用
func hello(ctx context.Context, client idl.HelloServiceClient) string {
	resp, err := client.SayHello(ctx, &idl.HelloRequest{})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return ""
	}
	return resp.Greeting
}

func main2() {
	InitGrpcJaeger()
	defer closer.Close()
	// 连接到GRPC服务端
	conn, err := grpc.Dial("127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(common.ClientTraceInterceptor),
	)
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}
	defer conn.Close()

	client := idl.NewHelloServiceClient(conn)

	userId := 8
	ctx := context.WithValue(context.Background(), "user_id", strconv.Itoa(userId))                       //"user_id"只在本进程内使用
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{"organization": "daqiaoqiao"})) //"organization"打算传给RPC的对端。注意通过grpc传递的metadata必须满足正则[0-9a-z-_.]+
	fmt.Println(hello(ctx, client))
}

// go run .\micro_service\trace\open_trace\rpc_tracer\client\
