package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// 实现credentials包里的PerRPCCredentials接口，即实现GetRequestMetadata和RequireTransportSecurity
type MyCredential struct {
	token string
}

func (c *MyCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": c.token,
	}, nil
}

// 是否需要TLS
func (*MyCredential) RequireTransportSecurity() bool {
	return true // 也可以为false，这样代码里就不需要写TLS相关代码了
}

func main() {
	// TLS认证
	creds, err := credentials.NewClientTLSFromFile(util.RootPath+"config/keys/server.crt", "")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	// 连接到GRPC服务端
	conn, err := grpc.Dial("localhost:5678",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}

	client := idl.NewHelloServiceClient(conn)
	resp, err := client.Login(context.Background(), &idl.LoginRequest{Name: "高性能golang", Pass: "123456"})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return
	}
	fmt.Println("登录成功")
	token := resp.Token
	conn.Close() //关闭第一个连接

	myCredential := &MyCredential{token: token}
	//重新创建一个新连接
	conn2, err := grpc.Dial("localhost:5678",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(myCredential), //附加安全信息的通用接口
	)
	client2 := idl.NewHelloServiceClient(conn2)
	resp2, err := client2.SayHello(context.Background(), &idl.HelloRequest{})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return
	}
	fmt.Println(resp2.Greeting)
	conn2.Close() //关闭第二个连接
}

// go run .\micro_service\secure\client\
