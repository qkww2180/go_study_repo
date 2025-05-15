package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	student_service "six/grpc/idl/my_proto"

	"google.golang.org/grpc"
)

type StudentServer struct {
	student_service.UnimplementedStudentServiceServer
}

func (s *StudentServer) GetStudentInfo(ctx context.Context, request *student_service.Request) (*student_service.Student, error) {
	defer func() {
		if err := recover(); err != nil { //避免子协程panic导致整个GRPC服务挂掉
			fmt.Printf("执行接口函数时出错: %v\n", err)
		}
	}()
	studentId := request.StudentId
	if len(studentId) == 0 {
		return nil, errors.New("studentId is empty") //参数错误，给调用方返回一个error
	}
	student := &student_service.Student{
		Name:   "张三",
		Age:    28,
		Height: 1.75,
	}
	return student, nil
}

func main() {
	// 监听本地的2346端口
	lis, err := net.Listen("tcp", "127.0.0.1:2346")
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	// 注册服务
	student_service.RegisterStudentServiceServer(server, new(StudentServer))
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}

// go run ./i_grpc
