package concurrence

import (
	"fmt"
	"time"
)

var asyncChann = make(chan int, 2) //缓冲长度为2，put可以比take多两次

func CloseChannel() {
	asyncChann <- 1
	go func() {
		time.Sleep(3 * time.Second)
		asyncChann <- 2
	}()
	// close(asyncChann) //channel关闭后读操作会立即返回，即使channel为空
	v, ok := <-asyncChann
	fmt.Println(v, ok)
	v, ok = <-asyncChann
	fmt.Println(v, ok)
	// asyncChann <- 3 //channel关闭后就不允许再写入元素了
}

func ChannelBlock() {
	go func() {
		time.Sleep(time.Second)
		if v, ok := <-asyncChann; ok {
			fmt.Println("take", v)
		}
		time.Sleep(time.Second)
		if v, ok := <-asyncChann; ok {
			fmt.Println("take", v)
		}
	}()
	asyncChann <- 1
	fmt.Println("send 1")
	asyncChann <- 2
	fmt.Println("send 2")
	asyncChann <- 3
	fmt.Println("send 3")
	time.Sleep(2 * time.Second)
}

func TraverseChannel() {
	asyncChann <- 1
	asyncChann <- 2
	close(asyncChann)

	for {
		if v, ok := <-asyncChann; ok {
			fmt.Println(v)
		} else { // channel已空，且已closed
			break
		}
	}
	fmt.Println("for finish")

	// go func() {
	// 	for {
	// 		asyncChann <- 10
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	// for ele := range asyncChann { //取走channel里的元素
	// 	fmt.Println(ele)
	// }
}
