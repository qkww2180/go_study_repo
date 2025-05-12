package test

import (
	gogo_pb "dqq/micro_service/pb/gogofaster/idl"
	std_pb "dqq/micro_service/pb/std/idl"
	"testing"
	"time"

	"github.com/cloudwego/fastpb"
	gogo_proto "github.com/gogo/protobuf/proto"
	google_proto "google.golang.org/protobuf/proto"
)

func BenchmarkGooglePb(b *testing.B) {
	student := std_pb.Student{
		Name:      "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, _ := google_proto.Marshal(&student) //序列化
		google_proto.Unmarshal(buf, &student)    //反序列化
	}
}

func BenchmarkGogoPb(b *testing.B) {
	student := gogo_pb.Student{
		Name:      "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, _ := gogo_proto.Marshal(&student) //序列化
		gogo_proto.Unmarshal(buf, &student)    //反序列化
	}
}

func BenchmarkFastPb(b *testing.B) {
	student := std_pb.Student{
		Name:      "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
	}
	b.ResetTimer()
	buf := make([]byte, student.Size())
	for i := 0; i < b.N; i++ {
		student.FastWrite(buf) //序列化
		student.Reset()
		fastpb.ReadMessage(buf, int8(fastpb.SkipTypeCheck), &student) //反序列化
	}
}

func TestFastPb(t *testing.T) {
	student := &std_pb.Student{
		Name:      "大乔乔",
		CreatedAt: time.Now().Unix(),
		Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
		Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
		Gender:    true,
		Age:       18,
		Height:    18.,
	}
	buf := make([]byte, student.Size()) //buf可以反复使用
	for i := 0; i < 10; i++ {
		offset := student.FastWrite(buf)                                                //序列化
		student.Reset()                                                                 //反序列化之前，先把fastpb.Reader清空
		_, err := fastpb.ReadMessage(buf[:offset], int8(fastpb.SkipTypeCheck), student) //反序列化
		if err != nil {
			t.Fatalf("%s", err)
		}
	}
}

// go test -v ./micro_service/pb/test -run=TestFastPb -count=1
// go test ./micro_service/pb/test -bench=^Benchmark.Pb$ -run=^$ -count=1 -benchmem -benchtime=2s
