package json

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/sonic"
)

var user = User{Name: "大乔乔", Age: 18, height: 170.5, Birthday: time.Now(), CreatedAt: MyDate(time.Now())}

func TestStdJson(t *testing.T) {
	bs, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("序列化失败:%s\n", err)
		return
	}
	str := string(bs)
	fmt.Println("序列化成功", str)

	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		fmt.Printf("反序列化失败:%s\n", err)
	} else {
		fmt.Println("反序列化成功", user) //使用Println直接打印user时会调用它的String()方法，user的String()又会调用各个成员变量的String()
	}
}

func TestSonicJson(t *testing.T) {
	bs, err := sonic.Marshal(user) //sonic的使用方法跟标准库完全一样，只需要替换包名
	if err != nil {
		fmt.Printf("序列化失败:%s\n", err)
		return
	}
	str := string(bs)
	fmt.Println("序列化成功", str)

	err = sonic.Unmarshal([]byte(str), &user)
	if err != nil {
		fmt.Printf("反序列化失败:%s\n", err)
	} else {
		fmt.Println("反序列化成功", user) //使用Println直接打印user时会调用它的String()方法，user的String()又会调用各个成员变量的String()
	}
}

func BenchmarkStdJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := json.Marshal(user)
		var inst User
		json.Unmarshal(bs, &inst)
	}
}

func BenchmarkSonicJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, _ := sonic.Marshal(user)
		var inst User
		sonic.Unmarshal(bs, &inst)
	}
}

// go test -v ./io/json/ -run=^TestStdJson$ -count=1
// go test -v ./io/json/ -run=^TestSonicJson$ -count=1
// go test ./io/json/ -bench=^Benchmark.*Json$ -run=^$ -count=1 -benchmem -benchtime=3s
/**
BenchmarkStdJson-8        520320              2278 ns/op
BenchmarkSonicJson-8      904772              1401 ns/op
*/
