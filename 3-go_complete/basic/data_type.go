package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func serialize1(arg any) ([]byte, error) {
	switch v := arg.(type) {
	case int:
		buffer := bytes.NewBuffer(make([]byte, 0, 8))
		if err := binary.Write(buffer, binary.BigEndian, int64(v)); err != nil {
			return nil, err
		} else {
			return buffer.Bytes(), nil
		}
	case string:
		return []byte(v), nil
	default:
		return nil, errors.New("不支持的数据类型")
	}
}

func serialize2(arg any) ([]byte, error) {
	v := reflect.ValueOf(arg)
	switch v.Kind() {
	case reflect.Int:
		buffer := bytes.NewBuffer(make([]byte, 0, 8))
		if err := binary.Write(buffer, binary.BigEndian, int64(v.Interface().(int))); err != nil {
			return nil, err
		} else {
			return buffer.Bytes(), nil
		}
	case reflect.String:
		return []byte(v.Interface().(string)), nil
	default:
		return nil, errors.New("不支持的数据类型")
	}
}

// 泛型
type IntString interface {
	int | string
}

func serialize3[T IntString](arg T) ([]byte, error) {
	return serialize1(arg)
	// return serialize2(arg)
}

type Number interface {
	int | float32
}

// 使用泛型的典型场景是不需要判断参数的具体类型
func add[T Number](arg1, arg2 T) T {
	return arg1 + arg2
}

func main13() {
	i := 8
	s := "A"
	fmt.Println(serialize1(i))
	fmt.Println(serialize1(s))
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println(serialize2(i))
	fmt.Println(serialize2(s))
	fmt.Println(strings.Repeat("-", 50))

	fmt.Println(serialize3(i))
	fmt.Println(serialize3(s))
	fmt.Println(strings.Repeat("-", 50))

	add(8, 9)
}
