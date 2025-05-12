package test

import (
	"fmt"
	"my_kitex/kitex_gen/math_service"
	"testing"

	"github.com/cloudwego/fastpb"
	"google.golang.org/protobuf/proto"
)

func TestGooglePb(t *testing.T) {
	req := math_service.AddRequest{Left: 5, Right: 3}
	if buf, err := proto.Marshal(&req); err == nil { //序列化
		if err = proto.Unmarshal(buf, &req); err != nil { //反序列化
			t.Errorf("google pb反序列化失败:%s", err)
		} else {
			fmt.Println(req.Left, req.Right)
		}
	} else {
		t.Errorf("google pb序列化失败:%s", err)
	}
}

func TestFastPb(t *testing.T) {
	req := math_service.AddRequest{Left: 5, Right: 3}
	buf := make([]byte, req.Size())
	offset := req.FastWrite(buf)                                                 //序列化
	req.Reset()                                                                  //如果要复用req，反序列化之前先把req清空
	_, err := fastpb.ReadMessage(buf[:offset], int8(fastpb.SkipTypeCheck), &req) //反序列化
	if err != nil {
		t.Errorf("fastpb反序列化失败:%s", err)
	} else {
		fmt.Println(req.Left, req.Right)
	}
}

func BenchmarkGooglePb(b *testing.B) {
	req := math_service.AddRequest{Left: 5, Right: 3}
	for i := 0; i < b.N; i++ {
		buf, _ := proto.Marshal(&req) //序列化
		proto.Unmarshal(buf, &req)    //反序列化
	}
}

func BenchmarkFastPb(b *testing.B) {
	req := math_service.AddRequest{Left: 5, Right: 3}
	buf := make([]byte, req.Size())
	for i := 0; i < b.N; i++ {
		offset := req.FastWrite(buf)                                       //序列化
		fastpb.ReadMessage(buf[:offset], int8(fastpb.SkipTypeCheck), &req) //反序列化
	}
}

// go test -v ./kitex_gen/math_service/test -run=Test.Pb -count=1
// go test -v ./kitex_gen/math_service/test -bench=Benchmark.Pb -run=^$ -count=1 -benchmem -benchtime=2s
/*
BenchmarkGooglePb-8      9377504               151.4 ns/op
BenchmarkFastPb-8       33753184                32.36 ns/op
*/
