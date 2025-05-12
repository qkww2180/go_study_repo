package io

import (
	"encoding/json"
	"fmt"
	"time"
)

//自定义json里的时间格式，核心是自定义一个type，实现MarshalJSON和UnmarshalJSON这两个方法

var MyDateFormat = "2006-01-02"

type MyDate time.Time

func (d MyDate) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", time.Time(d).Format(MyDateFormat))
	return []byte(s), nil
}

// 要改变自己，必须传指针
func (d *MyDate) UnmarshalJSON(bs []byte) (err error) {
	now, err := time.ParseInLocation(`"`+MyDateFormat+`"`, string(bs), time.Local) //注意MyDateFormat前后还得加引号
	*d = MyDate(now)                                                               // 要改变自己
	return
}

// 在print(MyDate)时会调用String()方法。print(User)时会间接地调到MyDate的String()方法
func (d MyDate) String() string {
	return time.Time(d).Format(MyDateFormat)
}

// 序列化就是把变量（包括基础类型的变量和复杂的结构体变量）转为二进制流（二进制流可以跟[]byte等价，[]byte又可以跟string等价），以便写入磁盘或发送到网络上。反序列化跟该过程相反。

type User struct {
	Name      string    //默认的json字段名跟原始变量名保持一致
	Age       int       `json:"-"`
	height    float32   //不可导出成员不会被序列化（该变量的值不会被导出到磁盘或网络上），否则就违背了“不可导出”的本意
	Birthday  time.Time //格式： 2023-09-29T20:14:11.7074482+08:00
	CreatedAt MyDate    //格式： 2023-09-29
	Sex       int       `json:"gender"`
	// Ch        chan int  //chan类型不支持序列化，Marshal()会返回error
	Address struct { //匿名结构体
		Province string
		City     string
	}
}

// 标准库json序列化背后使用的核心技术是反射。通过反射可以在运行时动态获得结构体成员变量的名称、(json)tag、是否可导出，可以获取成员变量的值，还可以调用结构体的方法(比如MarshalJSON和UnmarshalJSON)

func JsonSerialize() {
	var user = User{Name: "大乔乔", Age: 18, height: 170.5, Birthday: time.Now(), CreatedAt: MyDate(time.Now()), Sex: 1, Address: struct {
		Province string
		City     string
	}{Province: "河南", City: "郑州"}}

	bs, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("序列化失败:%s\n", err)
		return
	}
	fmt.Println("序列化成功", string(bs))

	var u User
	err = json.Unmarshal(bs, &u)
	if err != nil {
		fmt.Printf("反序列化失败:%s\n", err)
	} else {
		fmt.Printf("反序列化成功 %+v\n", u) //使用Println直接打印user时会调用它的String()方法，user的String()又会调用各个成员变量的String()
	}
}
