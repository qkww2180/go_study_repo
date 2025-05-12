package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found error") //error变量名一般以Err开头
	ErrServer   = errors.New("server error")
)

// 自定义Error
type MyError struct {
	Name string
	Code int
	Desc string
}

// 构造函数
func NewMyError(name string, code int, desc string) *MyError {
	return &MyError{
		Name: name,
		Code: code,
		Desc: desc,
	}

	// err:=new(NewMyError)
	// err.Name=name
	// err.Code=code
	// err.Desc=desc
	// return err
}

// 拥有Error() string即实现了error接口。Print error时默认会调用error的Error()方法
func (e MyError) Error() string {
	return fmt.Sprintf("[%d]%s: %s", e.Code, e.Name, e.Desc)
}

// Print对象时默认会调用对象的String()方法。error是个例外
func (e MyError) String() string {
	return e.Name
}

// 函数有多个返回值时，error通常是最后一个
func divide(a, b int) (int, error) {
	if b == 0 {
		// return 0, &MyError{}
		return 0, NewMyError("math", 101, "divide by zero")
		// return 0, ErrNotFound
		// return 0, ErrServer
		// return 0, fmt.Errorf("divide error %d %d", a, b)
	} else {
		return a / b, nil
	}
}

func main20() {
	fmt.Println(NewMyError("math", 101, "divide by zero"))

	c, err := divide(5, 0)
	if err != nil {
		fmt.Printf("出错 %s\n", err) //默认会调用err.Error()
		fmt.Printf("出错 %s\n", err.Error())
	} else {
		fmt.Printf("结果 %d\n", c)
	}
}
