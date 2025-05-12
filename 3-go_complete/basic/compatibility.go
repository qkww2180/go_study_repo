package main

//接口兼容，用interface{}还是泛型？

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// 服务方
func service(id any) (int, error) {
	var idInt int
	switch v := id.(type) {
	case int:
		idInt = v
	case string:
		if a, err := strconv.Atoi(v); err != nil {
			return 0, err
		} else {
			idInt = a
		}
	default:
		return 0, fmt.Errorf("非法的数据类型:%s", reflect.TypeOf(id).Name())
	}
	return idInt * 2, nil
}

// 调用方
func call[T int | string | bool](id T) {
	if r, err := service(id); err != nil {
		log.Printf("调服务接口失败:%s", err)
	} else {
		log.Printf("服务接口返回%d", r)
	}
}

// func main() {
// 	call("4")
// 	call(4)
// 	// call(4.0)
// 	call(false)
// }
