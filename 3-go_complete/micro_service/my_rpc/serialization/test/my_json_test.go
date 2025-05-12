package test

import (
	"dqq/micro_service/my_rpc/serialization"
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
	Sex  byte `json:"gender"`
}

type Book struct {
	ISBN     string `json:"isbn"`
	Name     string
	Price    float32  `json:"price"`
	Author   *User    `json:"author"` //把指针去掉试试
	Keywords []string `json:"kws"`
	Local    map[int]bool
}

func TestMyJson(t *testing.T) {
	user := User{
		Name: "钱钟书",
		Age:  57,
		Sex:  1,
	}
	book := Book{
		ISBN:     "4243547567",
		Name:     "围城",
		Price:    34.8,
		Author:   &user,                      //改成nil试试
		Keywords: []string{"爱情", "民国", "留学"}, //把这一行注释掉试一下，测测null
		Local:    map[int]bool{2: true, 3: false},
	}

	s := new(serialization.MyJson)

	if bytes, err := s.Marshal(user); err != nil { //也可以给Marshal传指针类型
		fmt.Printf("序列化失败: %v\n", err)
	} else {
		fmt.Println(string(bytes))
		var u User
		if err = s.Unmarshal(bytes, &u); err != nil {
			fmt.Printf("反序列化失败: %v\n", err)
		} else {
			fmt.Printf("user name %s\n", u.Name)
		}
	}

	if bytes, err := json.Marshal(user); err != nil {
		fmt.Printf("序列化失败: %v\n", err)
	} else {
		fmt.Println(string(bytes))
		var u User
		if err = json.Unmarshal(bytes, &u); err != nil {
			fmt.Printf("反序列化失败: %v\n", err)
		} else {
			fmt.Printf("user name %s\n", u.Name)
		}
	}

	if bytes, err := s.Marshal(book); err != nil { //也可以给Marshal传指针类型
		fmt.Printf("序列化失败: %v\n", err)
	} else {
		fmt.Println(string(bytes))
		var b Book //必须先声明值类型，再通过&给Unmarshal传一个指针参数。因为声明值类型会初始化为0值，而声明指针都没有创建底层的内存空间
		if err = s.Unmarshal(bytes, &b); err != nil {
			fmt.Printf("反序列化失败: %v\n", err)
		} else {
			fmt.Printf("book name %s author name %s local %v\n", b.Name, b.Author.Name, b.Local)
		}
	}

	if bytes, err := json.Marshal(book); err != nil {
		fmt.Printf("序列化失败: %v\n", err)
	} else {
		fmt.Println(string(bytes))
		var b Book
		if err = json.Unmarshal(bytes, &b); err != nil {
			fmt.Printf("反序列化失败: %v\n", err)
		} else {
			fmt.Printf("book name %s author name %s local %v\n", b.Name, b.Author.Name, b.Local)
		}
	}
}
