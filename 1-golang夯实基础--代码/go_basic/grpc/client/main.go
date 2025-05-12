package main

import (
	"context"
	grpc_service "dqq/go/basic/grpc/idl/service"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func timer(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	begin := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf("use time %d ms\n", time.Since(begin).Milliseconds())
	return err
}

func counter(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Printf("%s add once\n", method)
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

// client必须传正确的开发者key，才能调接口
func devKey(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "dev_key", "123456")
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("data/server.crt", "")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	//连接到服务端
	conn, err := grpc.NewClient(
		"localhost:5678", //根证书里的域名保持一致
		grpc.WithTransportCredentials(insecure.NewCredentials()), //Credential即使为空，也必须设置
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024), //默认情况下SendMsg上限是MaxInt32，RecvMsg上限是4M
			grpc.MaxCallSendMsgSize(1024),
		),
		// grpc.WithUnaryInterceptor(timer),
		grpc.WithChainUnaryInterceptor(timer, counter, devKey), //链式拦截器
		grpc.WithTransportCredentials(creds),                   //TLS数据加密
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}

	// const C = 3
	// wg := sync.WaitGroup{}
	// wg.Add(C)
	// for i := 0; i < C; i++ {
	// 	go func() {
	// 		defer wg.Done()
	//一个*grpc.ClientConn可以被多个Client公用--多路复用
	client := grpc_service.NewStudentClient(conn)
	//准备context
	ctx := context.Background()
	ctx1, cancel := context.WithTimeout(ctx, 1000*time.Millisecond) //调用超时
	defer cancel()
	//准备request
	request := grpc_service.QueryStudentRequest{
		Id:   123,
		Name: "dqq",
	}
	resp, err := client.QueryStudent(ctx1, &request, grpc.MaxCallRecvMsgSize(1024))
	if err != nil {
		fmt.Printf("rpc error: %s\n", err)
	} else {
		fmt.Printf("response: %+v\n", resp)
	}
	// 	}()
	// }
	// wg.Wait()

	streaming(client)
}

func streaming(client grpc_service.StudentClient) {
	request2 := grpc_service.StudentIds{Ids: []int64{100, 300, 500, 700, 900, 1000}}
	//流式地接收response
	stream2, err := client.QueryStudents2(context.Background(), &request2)
	if err != nil {
		fmt.Printf("build stream2 failed: %s", err)
	} else {
		for {
			stu, err := stream2.Recv() //从响应流中取得一个结果
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("recv response failed: %s\n", err)
				continue
			}
			fmt.Printf("response: %+v\n", stu)
		}
	}
	fmt.Println()

	//流式地发送request
	stream3, err := client.QueryStudents3(context.Background())
	if err != nil {
		fmt.Printf("build stream3 failed: %s", err)
	} else {
		for i := 1; i < 5; i++ {
			request := grpc_service.StudentId{Id: 123}
			stream3.Send(&request)
		}
		resp, err := stream3.CloseAndRecv() //关闭流，然后等待Server一次性返回全部结果
		if err != nil {
			fmt.Printf("recv response failed: %s", err)
		} else {
			for _, response := range resp.Students {
				fmt.Println(response.Id)
			}
		}
	}
	fmt.Println()

	//流式地发送request, 流式地接收response
	stream4, err := client.QueryStudents4(context.Background())
	done := make(chan struct{})
	if err != nil {
		fmt.Printf("build stream4 failed: %s", err)
	} else {
		go func() { //发送和接收同时进行
			for {
				resp, err := stream4.Recv() //从响应流中取得一个结果
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Printf("recv response failed: %s\n", err)
					continue
				}
				fmt.Println(resp.Id)
			}
			done <- struct{}{} //取出所有结果后对外发送一个信号
		}()
		for i := 1; i < 5; i++ {
			request := grpc_service.StudentId{Id: 123}
			stream4.Send(&request)
		}
	}
	stream4.CloseSend() //关闭流。客户端创建的stream最终都要调close，server端不用调close。正如http client需要close response body，而http server不需要close request body
	<-done              //done之后往下走
	fmt.Println()
}

// go run ./grpc/client
