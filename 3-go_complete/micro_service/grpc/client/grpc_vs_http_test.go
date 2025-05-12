package main

import (
	"bytes"
	"context"
	student_service "dqq/micro_service/grpc"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const student_id = 10

func TestGrpc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建client
	client := student_service.NewStudentServiceClient(conn)

	//准备好请求参数
	request := student_service.StudentID{Id: student_id}
	//发送请求，取得响应
	response, err := client.GetStudent(context.Background(), &request)
	if err != nil {
		fmt.Printf("get student failed: %s", err)
	} else {
		fmt.Println(response.Id)
	}
}

func TestHttp(t *testing.T) {
	client := http.Client{}

	//准备好请求参数
	sid := student_service.StudentID{Id: student_id}
	bs, err := sonic.Marshal(&sid) //也可以使用pb进行序列化，然后Content-type指定为"application/protobuf"
	if err != nil {
		fmt.Printf("marshal request failed: %s\n", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5679/", bytes.NewReader(bs))
	if err != nil {
		fmt.Printf("build request failed: %s\n", err)
		return
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("http rpc failed: %s\n", err)
		return
	}
	bs, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response stream failed: %s\n", err)
		return
	}
	resp.Body.Close()

	var stu student_service.Student
	err = sonic.Unmarshal(bs, &stu)
	if err != nil {
		fmt.Printf("unmarshal student failed: %s\n", err)
		return
	}
	fmt.Println(stu.Id)
}

const P = 10

func BenchmarkGrpc(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond) //连接超时设置为1000毫秒
	defer cancel()
	//连接到服务端
	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), //grpc.WithBlock()直到连接真正建立才会返回，否则连接是异步建立的。因此grpc.WithBlock()和Timeout结合使用才有意义
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建client
	client := student_service.NewStudentServiceClient(conn)
	for i := 0; i < b.N; i++ {
		const P = 10
		wg := sync.WaitGroup{}
		wg.Add(P)
		for j := 0; j < P; j++ {
			go func() {
				defer wg.Done()
				//准备好请求参数
				request := student_service.StudentID{Id: student_id}
				//调用服务，取得结果，结果反序列化为结构体
				client.GetStudent(context.Background(), &request)
			}()
		}
		wg.Wait()
	}
}

func BenchmarkHttp(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(P)
		for j := 0; j < P; j++ {
			go func() {
				defer wg.Done()
				client := http.Client{}
				//准备好请求参数
				sid := student_service.StudentID{Id: student_id}
				bs, err := sonic.Marshal(&sid)
				if err != nil {
					fmt.Printf("marshal request failed: %s\n", err)
					return
				}

				request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5679/", bytes.NewReader(bs))
				if err != nil {
					fmt.Printf("build request failed: %s\n", err)
					return
				}
				//调用服务
				resp, err := client.Do(request)
				if err != nil {
					fmt.Printf("http rpc failed: %s\n", err)
					return
				}
				//取得结果
				bs, err = io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("read response stream failed: %s\n", err)
					return
				}
				resp.Body.Close()
				//结果反序列化为结构体
				var stu student_service.Student
				sonic.Unmarshal(bs, &stu)
			}()
		}
		wg.Wait()
	}
}

// go test -v .\micro_service\grpc\client\ -run=^TestGrpc$ -count=1
// go test -v .\micro_service\grpc\client\ -run=^TestHttp$ -count=1
// go test .\micro_service\grpc\client\ -bench=^BenchmarkGrpc$ -run=^$ -count=1 -benchmem -benchtime=5s
// go test .\micro_service\grpc\client\ -bench=^BenchmarkHttp$ -run=^$ -count=1 -benchmem -benchtime=5s
// go test .\micro_service\grpc\client\ -bench=^Benchmark -run=^$ -count=1 -benchmem -benchtime=5s

/**
10个并发
BenchmarkGrpc-8            10000            527767 ns/op           56212 B/op       1101 allocs/op
BenchmarkHttp-8             3106          12205980 ns/op          121910 B/op        944 allocs/op
1个并发
BenchmarkGrpc-8            10000            518332 ns/op           56215 B/op       1102 allocs/op
BenchmarkHttp-8            69291             99370 ns/op            7439 B/op         63 allocs/op
*/
