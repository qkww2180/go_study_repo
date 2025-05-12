package grpcconnpool_test

import (
	"context"
	"fmt"
	rand "math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	student_service "dqq/micro_service/grpc"
	grpcconnpool "dqq/micro_service/grpc_conn_pool"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	pool *grpcconnpool.GrpcClientPool
)

func init() {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(10<<20), grpc.MaxCallRecvMsgSize(10<<20)),
	}
	cfg := &grpcconnpool.ClientConfig{}
	cfg.WithServerAddr("127.0.0.1:5678")
	cfg.WithConnectionPoolSize(10)
	cfg.WithConnectionLifeTime(10 * time.Minute)
	cfg.WithConnectionLifeTimeDeviation(60 * time.Second)
	pool = grpcconnpool.NewGrpcClientPool(*cfg, opts, &grpcconnpool.RoundRobin{})
}

func TestGrpcConnPool(t *testing.T) {
	conn := pool.GetConn()
	if conn == nil {
		t.Error("client is nil")
	}
	client := student_service.NewStudentServiceClient(conn)
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

// 并发测试
func TestGrpcConnPoolConncurent(t *testing.T) {
	const P = 1000
	const LOOP = 1000
	var fail int32
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			conn := pool.GetConn()
			if conn == nil {
				t.Error("client is nil")
			}
			client := student_service.NewStudentServiceClient(conn)
			for j := 0; j < LOOP; j++ {
				//准备好请求参数
				request := student_service.StudentID{Id: 888}
				//调用服务，取得结果，结果反序列化为结构体
				_, err := client.GetStudent(context.Background(), &request)
				if err != nil {
					conn = pool.GetConn()
					if conn == nil {
						t.Error("client is nil")
					}
					client = student_service.NewStudentServiceClient(conn)
					_, err := client.GetStudent(context.Background(), &request) //允许失败一次，因获取到连接后，可能连接又被关闭了
					if err != nil {
						// t.Errorf("get student failed: %s", err)
						atomic.AddInt32(&fail, 1)
					}
				}
				if rand.Float32() < 0.01 { //故意制造异常，让连接不可用
					conn.Close()
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("fail %d ratio %f", fail, float64(fail)/(P*LOOP))
}

// go test -v ./micro_service/grpc_conn_pool/ -run=^TestGrpcConnPool$ -count=1
// go test -v ./micro_service/grpc_conn_pool/ -run=^TestGrpcConnPoolConncurent$ -count=1
