package main

import (
	"errors"
	"fmt"
	"strconv"
)

// 演示%w和errors.Is()
func errorWrap() {
	a := errors.New("大乔乔")
	b := fmt.Errorf("%w", a)
	c := fmt.Errorf("123 %w abc", b)
	fmt.Printf("%t\n", errors.Is(c, a)) //true
	fmt.Printf("%t\n", errors.Is(c, b)) //true
}

func a(s string) (int, error) {
	result, err := b(s)
	if err != nil {
		// return 0, err
		return 0, fmt.Errorf("a-> (%w)", err)
	}
	return result, nil
}

func b(s string) (int, error) {
	result, err := c(s)
	if err != nil {
		// return 0, err
		return 0, fmt.Errorf("b-> (%w) (%w)", err, ErrNotFound)
	}
	return result, nil
}

func c(s string) (int, error) {
	result, err := strconv.Atoi(s)
	if err != nil {
		// return 0, err
		return 0, fmt.Errorf("c:cast-> (%w)", err) //有多个%w，则errors.Is()的target为任意一个时都返回True
	}

	if result == 0 {
		err = errors.New("divide by zero")
		return 0, fmt.Errorf("c:divide-> %s (%w)", err.Error(), ErrServer)
	}
	return 10 / result, nil
}

func handler() int {
	if result, err := a("0a3"); err != nil {
		fmt.Println(err) //准备处理error，即不打算把error再往上抛了，仅在此处打印error

		// if err == ErrNotFound {
		if errors.Is(err, ErrNotFound) {
			return 400
		}
		if errors.Is(err, ErrServer) { // %w
			return 500
		}

		var e MyError
		e, ok := err.(MyError) //类型断言
		if ok {
			fmt.Println(e.Name)
		}

		//想取得error的成员变量，用errors.As()
		var e2 MyError
		ok = errors.As(err, &e2) //errors.As背后其实就是类型断言
		if ok {
			fmt.Println(e2.Name)
		}

		return 200
	} else {
		return result
	}
}

func main29() {
	errorWrap()
	fmt.Println(handler())
}
