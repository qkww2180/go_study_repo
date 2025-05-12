package main

import (
	"context"
	student_service "dqq/micro_service/grpc"
	"dqq/micro_service/idl"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client  student_service.StudentServiceClient
	client2 idl.HelloServiceClient
)

/**
多路复用GRPC使用HTTP/2作为应用层的传输协议，HTTP/2会复用底层的TCP连接。每一次RPC调用会产生一个新的Stream，每个Stream包含多个Frame，Frame是HTTP/2里面最小的数据传输单位。同时每个Stream有唯一的ID标识，如果是客户端创建的则ID是奇数，服务端创建的ID则是偶数。如果一条连接上的ID使用完了，Client会新建一条连接，Server也会给Client发送一个GOAWAY Frame强制让Client新建一条连接。一条GRPC连接允许并发的发送和接收多个Stream，而控制的参数便是MaxConcurrentStreams，Golang的服务端默认是100。

超时重连我们在通过调用Dial或者DialContext函数创建连接时，默认只是返回ClientConn结构体指针，同时会启动一个Goroutine异步的去建立连接。如果想要等连接建立完再返回，可以指定grpc.WithBlock()传入Options来实现。超时机制很简单，在调用的时候传入一个timeout的context就可以了。重连机制通过启动一个Goroutine异步的去建立连接实现的，可以避免服务器因为连接空闲时间过长关闭连接、服务器重启等造成的客户端连接失效问题。也就是说通过GRPC的重连机制可以完美的解决连接池设计原则中的空闲连接的超时与保活问题。

GRPC参数调优：
MaxSendMsgSizeGRPC最大允许发送的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。
MaxRecvMsgSizeGRPC最大允许接收的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。
InitialWindowSize基于Stream的滑动窗口，类似于TCP的滑动窗口，用来做流控，默认64KiB，吞吐量上不去，Client和Server我们调到1GiB。
InitialConnWindowSize基于Connection的滑动窗口，默认16 * 64KiB，吞吐量上不去，Client和Server我们也都调到1GiB。
KeepAliveTime每隔KeepAliveTime时间，发送PING帧测量最小往返时间，确定空闲连接是否仍然有效，我们设置为10S。
KeepAliveTimeout超过KeepAliveTimeout，关闭连接，我们设置为3S。
PermitWithoutStream如果为true，当连接空闲时仍然发送PING帧监测，如果为false，则不发送忽略。我们设置为true。
*/

func InitClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()), //Credential即使为空，也必须设置
		grpc.WithBlock(), //grpc.WithBlock()直到连接真正建立才会返回，否则连接是异步建立的。因此grpc.WithBlock()和Timeout结合使用才有意义。server端正常的情况下使用grpc.WithBlock()得到的connection.GetState()为READY，不使用grpc.WithBlock()得到的connection.GetState()为IDEL
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(10<<20), grpc.MaxCallRecvMsgSize(10<<20)), //默认情况下SendMsg上限是MaxInt32，RecvMsg上限是4M，这里都修改为10M
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建client
	client = student_service.NewStudentServiceClient(conn)
	client2 = idl.NewHelloServiceClient(conn)
}

func general() {
	//准备好请求参数
	request := student_service.StudentID{Id: 10}
	//发送请求，取得响应
	response, err := client.GetStudent(context.Background(), &request)
	if err != nil {
		fmt.Printf("get student failed: %s", err)
	} else {
		fmt.Println(response.Id)
	}
	fmt.Println()
}

func exception() {
	request2 := student_service.StudentIDs{Ids: []int32{}} // 参数为空，故意触发服务端返回error
	response2, err := client.GetStudents(context.Background(), &request2)
	if err != nil {
		fmt.Printf("get students failed: %s", err)
	} else {
		for _, response := range response2.Data {
			fmt.Println(response.Id)
		}
	}
	fmt.Println()
}

func multiplexing() {
	//grpc基于http2协议，本身支持多路复用，N个goroutine用一个HTTP2连接没有任何问题，不会单纯因为没有可用连接而阻塞执行。
	const CLIENT_COUNT = 10 //10个client并发调用
	wg := sync.WaitGroup{}
	wg.Add(CLIENT_COUNT)
	for i := 0; i < CLIENT_COUNT; i++ {
		go func(i int) {
			defer wg.Done()
			request2 := student_service.StudentIDs{Ids: []int32{100, 300, 500, 700, 900, 1000}}
			fmt.Printf("client %d start\n", i)
			_, err := client.GetStudents(context.Background(), &request2) //client支持多协程并发使用
			if err != nil {
				fmt.Printf("get students failed: %s", err)
				return
			}
			fmt.Printf("client %d finish\n", i)
		}(i)
	}
	wg.Wait()
	fmt.Println()
}

func multiplexing2() {
	request := student_service.StudentID{Id: 10}
	//发送请求，取得响应
	if resp, err := client.GetStudent(context.Background(), &request); err == nil {
		fmt.Println(resp.Height)
	}
	request2 := idl.HelloRequest{}
	//发送请求，取得响应
	if resp2, err := client2.SayHello(context.Background(), &request2); err == nil {
		fmt.Println(resp2.Greeting)
	}
	fmt.Println()
}

func streaming() {
	request2 := student_service.StudentIDs{Ids: []int32{100, 300, 500, 700, 900, 1000}}
	//流式地接收response
	stream2, err := client.GetStudents2(context.Background(), &request2)
	if err != nil {
		fmt.Printf("build stream2 failed: %s", err)
	} else {
		for {
			resp, err := stream2.Recv() //从响应流中取得一个结果
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Printf("recv response failed: %s\n", err)
				continue
			}
			fmt.Println(resp.Id)
		}
	}
	fmt.Println()

	//流式地发送request
	stream3, err := client.GetStudents3(context.Background())
	if err != nil {
		fmt.Printf("build stream3 failed: %s", err)
	} else {
		for i := 1; i < 5; i++ {
			request := student_service.StudentID{Id: int32(i)}
			stream3.Send(&request)
		}
		resp, err := stream3.CloseAndRecv() //关闭流，然后等待Server一次性返回全部结果
		if err != nil {
			fmt.Printf("recv response failed: %s", err)
		} else {
			for _, response := range resp.Data {
				fmt.Println(response.Id)
			}
		}
	}
	fmt.Println()

	//流式地发送request, 流式地接收response
	stream4, err := client.GetStudents4(context.Background())
	done := make(chan struct{})
	if err != nil {
		fmt.Printf("build stream4 failed: %s", err)
	} else {
		go func() { //发送和接收同时进行
			for {
				resp, err := stream4.Recv() //从响应流中取得一个结果
				if err != nil {
					if err == io.EOF {
						done <- struct{}{} //取出所有结果后对外发送一个信号
						break
					}
					fmt.Printf("recv response failed: %s\n", err)
					continue
				}
				fmt.Println(resp.Id)
			}
		}()
		for i := 1; i < 5; i++ {
			request := student_service.StudentID{Id: int32(i)}
			stream4.Send(&request)
		}
	}
	stream4.CloseSend() //关闭流。客户端创建的stream最终都要调close，server端不用调close
	<-done              //done之后往下走
	fmt.Println()
}

func main1() {
	InitClient()
	general()
	exception()
	multiplexing()
	multiplexing2()
	streaming()
}

// go run ./micro_service/grpc/client
