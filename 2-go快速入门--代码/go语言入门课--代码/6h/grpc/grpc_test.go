package main

import (
	"context"
	"fmt"
	student_service "six/grpc/idl/my_proto"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestService(t *testing.T) {
	// 连接到GRPC服务端
	conn, err := grpc.NewClient("127.0.0.1:2346", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		t.Fail()
	}
	defer conn.Close()
	client := student_service.NewStudentServiceClient(conn)

	// 执行RPC调用并打印收到的响应数据，指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.GetStudentInfo(ctx, &student_service.Request{StudentId: "学生1"})
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		t.Fail()
	}
	fmt.Printf("Name %s Age %d Height %.1f\n", resp.Name, resp.Age, resp.Height)
}

// go test -v ./i_grpc
