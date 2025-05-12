package pb

import (
	fast_student_service "dqq/micro_service/pb/gogofaster/idl"
	std_student_service "dqq/micro_service/pb/std/idl"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	gogo_proto "github.com/gogo/protobuf/proto"
	google_proto "google.golang.org/protobuf/proto"
)

var student = std_student_service.Student{
	Name:      "大乔乔",
	CreatedAt: time.Now().Unix(),
	Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
	Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
	Gender:    true,
	Age:       18,
	Height:    18.,
}

var gogo_student = fast_student_service.Student{
	Name:      "大乔乔",
	CreatedAt: time.Now().Unix(),
	Locations: []string{"北京", "上海", "广东", "四川", "江苏"},
	Scores:    map[string]float32{"语文": 90., "英语": 80, "数学": 70},
	Gender:    true,
	Age:       18,
	Height:    18.,
}

func TestProtobuf(t *testing.T) {
	bs, err := google_proto.Marshal(&student) //是*Student实现了ProtoReflect()方法，所以这里要传指针
	if err != nil {
		fmt.Printf("proto序列化失败: %s", err)
		t.Fail()
	} else {
		fmt.Println(string(bs))
		var stu2 std_student_service.Student
		err = google_proto.Unmarshal(bs, &stu2)
		if err != nil {
			fmt.Printf("proto反序列化失败: %s", err)
			t.Fail()
		} else {
			fmt.Println(stu2.Locations)
		}
	}
}

func BenchmarkProtobuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := google_proto.Marshal(&student)
		var stu2 std_student_service.Student
		google_proto.Unmarshal(bs, &stu2)
	}
}

func BenchmarkJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := json.Marshal(&student)
		var stu2 std_student_service.Student
		json.Unmarshal(bs, &stu2)
	}
}

func BenchmarkSonicJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := sonic.Marshal(&student)
		var stu2 std_student_service.Student
		sonic.Unmarshal(bs, &stu2)
	}
}

func BenchmarkGogofasterProtobuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := gogo_proto.Marshal(&gogo_student) //等价于 bs, _ := gogo_student.Marshal()
		var stu2 fast_student_service.Student
		gogo_proto.Unmarshal(bs, &stu2)
	}
}

// go test -v .\micro_service\pb\ -run=^TestProtobuf$ -count=1
// go test .\micro_service\pb\ -bench=^BenchmarkProtobuf$ -run=^$ -benchmem -timeout=2s
// go test .\micro_service\pb\ -bench=^BenchmarkJson$ -run=^$ -benchmem
// go test .\micro_service\pb\ -bench=^BenchmarkSonicJson$ -run=^$ -benchmem
// go test .\micro_service\pb\ -bench=^BenchmarkGogofasterProtobuf$ -run=^$ -benchmem
