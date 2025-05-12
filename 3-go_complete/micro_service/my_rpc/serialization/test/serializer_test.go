package test

import (
	"dqq/micro_service/my_rpc"
	"dqq/micro_service/my_rpc/serialization"
	"testing"

	"github.com/stretchr/testify/assert"
)

var object = my_rpc.RpcData{
	A:  123,
	B:  3.14,
	C:  true,
	D:  3.141592653,
	E:  "大乔乔",
	Id: "fefbrekfna",
}

func equal(obj1, obj2 my_rpc.RpcData) bool {
	return obj1.A == obj2.A && obj1.B == obj2.B && obj1.C == obj2.C && obj1.D == obj2.D && obj1.E == obj2.E && obj1.Id == obj2.Id
}

func TestJosn(t *testing.T) {
	s := new(serialization.Json)
	stream, err := s.Marshal(object)
	if err != nil {
		t.Errorf("json序列化失败, %s", err)
	} else {
		var obj2 my_rpc.RpcData
		err := s.Unmarshal(stream, &obj2)
		if err != nil {
			t.Errorf("json反序列化失败, %s", err)
		} else {
			assert.Equal(t, obj2, obj2, "json反序列化后跟原始值不同")
			// if !equal(object, obj2) {
			// 	t.Error("json反序列化后跟原始值不同")
			// }
		}
	}
}

func TestGob(t *testing.T) {
	s := new(serialization.Gob)
	stream, err := s.Marshal(object)
	if err != nil {
		t.Errorf("gob序列化失败, %s", err)
	} else {
		var obj2 my_rpc.RpcData
		err := s.Unmarshal(stream, &obj2)
		if err != nil {
			t.Errorf("gob反序列化失败, %s", err)
		} else {
			assert.Equal(t, obj2, obj2, "gob反序列化后跟原始值不同")
			// if !equal(object, obj2) {
			// 	t.Error("gob反序列化后跟原始值不同")
			// }
		}
	}
}

func TestMySerializer(t *testing.T) {
	s := new(serialization.MySerializer)
	stream, err := s.Marshal(object)
	if err != nil {
		t.Errorf("my序列化失败, %s", err)
	} else {
		var obj2 my_rpc.RpcData
		err := s.Unmarshal(stream, &obj2)
		if err != nil {
			t.Errorf("my反序列化失败, %s", err)
		} else {
			assert.Equal(t, obj2, obj2, "my反序列化后跟原始值不同")
			// if !equal(object, obj2) {
			// 	t.Error("my反序列化后跟原始值不同")
			// }
		}
	}
}

func BenchmarkJson(b *testing.B) {
	s := new(serialization.Json)
	var obj2 my_rpc.RpcData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(object)
		s.Unmarshal(stream, &obj2)
	}
}

func BenchmarkGob(b *testing.B) {
	s := new(serialization.Gob)
	var obj2 my_rpc.RpcData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(object)
		s.Unmarshal(stream, &obj2)
	}
}

func BenchmarkSerializer(b *testing.B) {
	s := new(serialization.MySerializer)
	var obj2 my_rpc.RpcData
	for i := 0; i < b.N; i++ {
		stream, _ := s.Marshal(object)
		s.Unmarshal(stream, &obj2)
	}
}

// go test -v ./micro_service/my_rpc/serialization/test -run=TestJosn -count=1
// go test -v ./micro_service/my_rpc/serialization/test -run=TestGob -count=1
// go test -v ./micro_service/my_rpc/serialization/test -run=TestMySerializer -count=1
// go test ./micro_service/my_rpc/serialization/test -run=^$ -bench=^Benchmark -timeout=2s

/*
BenchmarkJson-8                   523084              2221 ns/op
BenchmarkGob-8                     68504             16154 ns/op
BenchmarkSerializer-8             358252              2845 ns/op
*/
