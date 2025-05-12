package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials" // 引入grpc认证包
	"google.golang.org/grpc/metadata"
)

var (
	defautHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

type MyServer struct{}

// 校验用户名和密码是否正确，如果正确则返回token
func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	username := request.Name
	// passwd := request.Pass
	//读数据库，校验用户名和密码是否正确
	header := defautHeader
	payload := JwtPayload{
		Issue:       "bilibili",
		IssueAt:     time.Now().Unix(),                         //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(), //3天后过期，需要重新登录
		UserDefined: map[string]any{"name": username, "role": "up主", "vip": true},
	}
	if token, err := GenJWT(header, payload); err != nil {
		log.Printf("生成token失败: %v", err)
		return nil, err
	} else {
		return &idl.LoginResponse{Token: token}, nil
	}
}

// 检查调用方是否登录成功，如何成功返回token。这里采用JWT校验，也可以采用更简单的token校验方式
func (*MyServer) checkIdentity(ctx context.Context) *JwtPayload {
	meta, ok := metadata.FromIncomingContext(ctx) //从调用方传过来的context里获取metadata
	if !ok {
		log.Println("从IncomingContext中获取metadata失败")
		return nil
	}
	value, ok := meta["token"] //跟调用方约定好key
	if !ok || len(value) == 0 {
		log.Println("metadata中没有token")
		return nil
	}
	token := value[0] //value是个切片，我们只取首元素

	_, payload, err := VerifyJwt(token)
	if err != nil {
		log.Printf("非法token %s\n", token)
		return nil
	}
	return payload
}

func (server *MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	payload := server.checkIdentity(ctx)
	if payload == nil {
		return nil, errors.New("身份认证失败")
	}
	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: fmt.Sprintf("你好 %s, 你的角色是 %s", payload.UserDefined["name"], payload.UserDefined["role"])}
	return resp, nil
}

func main() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "localhost:5678")
	if err != nil {
		panic(err)
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(util.RootPath+"config/keys/server.crt", util.RootPath+"config/keys/server.key")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(grpc.Creds(creds))
	// 注册服务
	idl.RegisterHelloServiceServer(server, new(MyServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run .\micro_service\secure\server\
