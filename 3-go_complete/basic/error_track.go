package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNotFound = errors.New("not found error")
	ErrServer   = errors.New("server error")
)

func a(s string) (int, error) {
	result, err := b(s)
	if err != nil {
		return 0, fmt.Errorf("a-> %w", err)
	}
	return result, nil
}

func b(s string) (int, error) {
	result, err := c(s)
	if err != nil {
		return 0, fmt.Errorf("b-> %w", err)
	}
	return result, nil
}

func c(s string) (int, error) {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("c:cast-> %w %w", err, ErrNotFound) //有多个%w，则errors.Is()的target为任意一个时都返回True
	}

	if result == 0 {
		err = errors.New("divide by zero")
		return 0, fmt.Errorf("c:divide-> %w %w", err, ErrServer)
	}

	return 10 / result, nil
}

func handler() int {
	if result, err := a("0"); err != nil {
		fmt.Println(err) //准备处理error，即不打算把error再往上抛了，仅在此处打印error

		if errors.Is(err, ErrNotFound) {
			return 400
		}
		if errors.Is(err, ErrServer) {
			return 500
		}
		return 200
	} else {
		return result
	}
}

func main7() {
	fmt.Println(handler())
}
