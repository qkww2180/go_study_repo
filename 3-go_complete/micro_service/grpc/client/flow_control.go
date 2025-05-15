package main

import (
	"context"
	"dqq/micro_service/idl"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DATA_SIZE = 1024 * 128
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(DATA_SIZE+4), grpc.MaxCallRecvMsgSize(DATA_SIZE+4)),
		//试验结论：InitialWindowSize和InitialConnWindowSize这俩参数不论对于Unary调用还是Stream调用都不起作用
		grpc.WithInitialWindowSize(DATA_SIZE*2),
		grpc.WithInitialConnWindowSize(DATA_SIZE*20),
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建client
	client := idl.NewEchoClient(conn)
	request := idl.EchoRequest{Data: string(make([]byte, DATA_SIZE))}
	resp, err := client.Hello(context.Background(), &request)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(len(resp.Data))
	}

	//流式地发送request
	stream, err := client.StreamHello(context.Background())
	if err != nil {
		fmt.Printf("build stream failed: %s", err)
	} else {
		for i := 1; i < 5; i++ {
			stream.Send(&request)
		}
		resp, err := stream.CloseAndRecv() //关闭流，然后等待Server一次性返回全部结果
		if err != nil {
			fmt.Printf("recv response failed: %s", err)
		} else {
			fmt.Println(len(resp.Data))
		}
	}

	time.Sleep(time.Second)
	begin := time.Now()
	const P = 1000
	const LOOP = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func(i int) {
			defer wg.Done()
			for j := 0; j < LOOP; j++ {
				stream, err := client.StreamHello(context.Background())
				if err != nil {
					fmt.Printf("build stream failed: %s", err)
				} else {
					for i := 1; i < 5; i++ {
						stream.Send(&request)
					}
					_, err := stream.CloseAndRecv() //关闭流，然后等待Server一次性返回全部结果
					if err != nil {
						fmt.Printf("recv response failed: %s", err)
					}
					// else {
					// 	fmt.Println(len(resp.Data))
					// }
				}
			}
		}(i)
	}
	wg.Wait()
	elapse := time.Since(begin).Milliseconds()
	fmt.Printf("use time %d ms, average %d call per ms\n", elapse, P*LOOP/elapse)
}

// go run ./micro_service/i_grpc/client
