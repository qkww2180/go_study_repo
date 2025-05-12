package main

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

type ETS struct{} //ETS跟标准库的struct{}等价（可以互换）

// 所有的空结构体指向同一个地址(内核是完全一样的)
func allEmptyStructIsSame() {
	var a ETS
	var b ETS
	var c struct{}
	fmt.Printf("address of a %p b %p c %p\n", &a, &b, &c)
	fmt.Printf("size of a %d b %d c %d\n", unsafe.Sizeof(a), unsafe.Sizeof(b), unsafe.Sizeof(c))
	fmt.Printf("size of a %d b %d c %d\n", reflect.TypeOf(a).Size(), reflect.TypeOf(b).Size(), reflect.TypeOf(c).Size())
}

// 空结构体的应用场景
func scenariosOfEmptyStruct() {
	set := map[int]struct{}{
		1: {},
		4: {},
		7: {},
	}
	if _, exists := set[5]; exists {
		fmt.Println("5是存在的")
	} else {
		fmt.Println("5是不存在的")
	}

	blocker := make(chan struct{})
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("done")
		blocker <- ETS{}
	}()
	<-blocker //等待子协程结束
}

func main26() {
	allEmptyStructIsSame()
	scenariosOfEmptyStruct()
}
