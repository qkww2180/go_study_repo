package grpc_model_test

import (
	grpc_model "dqq/go/basic/grpc/idl/model"
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

var (
	stu = &grpc_model.Student{
		Id:        1,
		Name:      "大乔乔",
		Locations: []string{"河南", "郑州"},
		Scores:    map[string]float32{"英文": 53.4},
	}
)

func TestStudentSerialize(t *testing.T) {
	// proto序列化
	bb, _ := proto.Marshal(stu)
	fmt.Println(string(bb))
	// stu1 := new(grpc_model.Student)
	var stu1 grpc_model.Student
	proto.Unmarshal(bb, &stu1)
	fmt.Println(stu1.Id, stu1.Name, stu1.Locations, stu1.Age)

	// protojson序列化
	bj, _ := protojson.Marshal(stu)
	fmt.Println(string(bj))
	stu2 := new(grpc_model.Student)
	protojson.Unmarshal(bj, stu2)
	fmt.Println(stu2.Id, stu2.Name, stu2.Locations, stu2.Age)

	// prototext序列化
	bp, _ := prototext.Marshal(stu)
	fmt.Println(string(bp))
	stu3 := new(grpc_model.Student)
	prototext.Unmarshal(bp, stu3)
	fmt.Println(stu3.Id, stu3.Name, stu3.Locations, stu3.Age)
}

func BenchmarkJson(b *testing.B) {
	buffer := new(grpc_model.Student)
	for i := 0; i < b.N; i++ {
		bs, _ := json.Marshal(stu)
		json.Unmarshal(bs, buffer)
	}
}

func BenchmarkProtoJson(b *testing.B) {
	buffer := new(grpc_model.Student)
	for i := 0; i < b.N; i++ {
		bs, _ := protojson.Marshal(stu)
		protojson.Unmarshal(bs, buffer)
	}
}

func BenchmarkProtoText(b *testing.B) {
	buffer := new(grpc_model.Student)
	for i := 0; i < b.N; i++ {
		bs, _ := prototext.Marshal(stu)
		prototext.Unmarshal(bs, buffer)
	}
}

func BenchmarkProto(b *testing.B) {
	buffer := new(grpc_model.Student)
	for i := 0; i < b.N; i++ {
		bs, _ := proto.Marshal(stu)
		proto.Unmarshal(bs, buffer)
	}
}

// go test -v ./grpc/idl/model -run=^TestStudentSerialize$ -count=1
// go test ./grpc/idl/model -bench=^Benchmark -run=^$ -count=1 -benchmem
/*
goos: windows
goarch: amd64
pkg: dqq/go/basic/grpc/idl/model
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkJson-8           454784              2340 ns/op             464 B/op         15 allocs/op
BenchmarkProtoJson-8      208869              5033 ns/op            1794 B/op         57 allocs/op
BenchmarkProtoText-8      200005              5478 ns/op            1953 B/op         68 allocs/op
BenchmarkProto-8          942536              1272 ns/op             432 B/op         17 allocs/op
*/
