package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"fmt"
	"log"
	"strconv"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func traceInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now()
	traceId := util.RandStringRunes(10) //生成随机字符串
	// 跨RPC传递的context
	ctx = metadata.AppendToOutgoingContext(ctx, "trace_id", traceId)
	ctx = metadata.AppendToOutgoingContext(ctx, "user_id", ctx.Value("user_id").(string))
	// 通过WithValue放入context里的数据不会通过RPC传播
	// ctx = context.WithValue(ctx, "trace_id", traceId)
	// ctx = context.WithValue(ctx, "user_id", ctx.Value("user_id").(string))

	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, ctx.Value("user_id").(string), method, begin.UnixNano(), time.Since(begin).Nanoseconds())
	return err
}

// RPC调用
func hello(ctx context.Context, client idl.HelloServiceClient) string {
	begin := time.Now()
	traceId := util.RandStringRunes(10) //生成随机字符串
	ctx = metadata.AppendToOutgoingContext(ctx, "trace_id", traceId)
	ctx = metadata.AppendToOutgoingContext(ctx, "user_id", ctx.Value("user_id").(string))
	defer func() {
		log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, ctx.Value("user_id").(string), "hello", begin.UnixNano(), time.Since(begin).Nanoseconds())
	}()
	resp, err := client.SayHello(ctx, &idl.HelloRequest{})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return ""
	}
	return resp.Greeting
}

func main() {
	// 连接到GRPC服务端
	conn, err := grpc.Dial("127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(traceInterceptor),
	)
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}
	defer conn.Close()

	client := idl.NewHelloServiceClient(conn)

	userId := 8
	ctx := context.WithValue(context.Background(), "user_id", strconv.Itoa(userId))
	fmt.Println(hello(ctx, client))
	client.Login(ctx, &idl.LoginRequest{})
}

// go run .\micro_service\trace\self_trace\rpc_tracer\client\
