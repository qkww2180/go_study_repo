package main

import (
	"context"
	student_service "dqq/micro_service/grpc"
	"dqq/micro_service/idl"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
)

type MyService struct {
	student_service.UnimplementedStudentServiceServer //StudentServiceServer接口要求实现mustEmbedUnimplementedStudentServiceServer()方法，而只有UnimplementedStudentServiceServer类实现了这个方法，所以必须继承它
}

func danger(b int32) {
	_ = 34 / b // 没有提防除0异常
}

func (s MyService) GetStudent(ctx context.Context, request *student_service.StudentID) (*student_service.Student, error) {
	defer func() {
		panicInfo := recover() //panicInfo是any类型，即传给panic()的参数
		if panicInfo != nil {
			fmt.Println(panicInfo)
		}
	}()
	danger(request.Id) //如果协程内发生panic，则下面的代码不会执行（导致返回的第一个参数为nil，客户端会接收到一个序列化错误），至于会不会Exit取决于有没有recover()
	// go danger(request.Id) //另一个协程内发生的panic，recover()捕获不了
	return &student_service.Student{Name: "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
		Id:        request.Id,
	}, nil
}

func (s MyService) GetStudents(ctx context.Context, request *student_service.StudentIDs) (*student_service.Students, error) {
	if len(request.Ids) == 0 {
		return nil, errors.New("please indicate Ids")
	}
	datas := make([]*student_service.Student, 0, len(request.Ids))
	for _, id := range request.Ids {
		student := &student_service.Student{Name: "大乔乔",
			CreatedAt: time.Now().Unix(),
			Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
			Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
			Gender:    true,
			Age:       18,
			Height:    18.,
			Id:        id,
		}
		datas = append(datas, student)
	}
	return &student_service.Students{Data: datas}, nil
}

func (s MyService) GetStudents2(request *student_service.StudentIDs, server student_service.StudentService_GetStudents2Server) error {
	for _, id := range request.Ids { //通过stream的形式返回多个结果
		student := &student_service.Student{Name: "大乔乔",
			CreatedAt: time.Now().Unix(),
			Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
			Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
			Gender:    true,
			Age:       18,
			Height:    18.,
			Id:        id,
		}
		err := server.Send(student) //向流中发送一个结果
		if err != nil {
			fmt.Printf("send Student %d failed: %s\n", id, err)
			return err
		}
	}
	return nil
}

func (s MyService) GetStudents3(server student_service.StudentService_GetStudents3Server) error {
	datas := make([]*student_service.Student, 0, 10)
	for {
		req, err := server.Recv() //从流中取出一个结果
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("recv request3 failed: %s\n", err)
			continue
		}
		student := &student_service.Student{Name: "大乔乔",
			CreatedAt: time.Now().Unix(),
			Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
			Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
			Gender:    true,
			Age:       18,
			Height:    18.,
			Id:        req.Id,
		}
		datas = append(datas, student)
	}
	return server.SendMsg(&student_service.Students{Data: datas})
}

func (s MyService) GetStudents4(server student_service.StudentService_GetStudents4Server) error {
	for {
		req, err := server.Recv() //从流中取出一个结果
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("recv request4 failed: %s\n", err)
			continue
		}
		student := &student_service.Student{Name: "大乔乔",
			CreatedAt: time.Now().Unix(),
			Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
			Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
			Gender:    true,
			Age:       18,
			Height:    18.,
			Id:        req.Id,
		}
		err = server.Send(student) //向流中发送一个结果
		if err != nil {
			fmt.Printf("send Student %d failed: %s\n", student.Id, err)
			return err
		}
	}
	return nil
}

func main() {
	// 监听本地的5678端口
	lis, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}
	//创建服务
	server := grpc.NewServer()
	// 注册服务的具体实现
	student_service.RegisterStudentServiceServer(server, &MyService{})
	// idl.RegisterHelloServiceServer(server, new(MyServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

type MyServer struct {
	idl.UnimplementedHelloHttpServer
}

func (*MyServer) Login(ctx context.Context, request *idl.LoginRequest) (*idl.LoginResponse, error) {
	return nil, nil
}

func (*MyServer) SayHello(ctx context.Context, request *idl.HelloRequest) (*idl.HelloResponse, error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	resp := &idl.HelloResponse{Greeting: "hello"}
	return resp, nil
}

// go run ./micro_service/grpc/server
