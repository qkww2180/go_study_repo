package concurrence

import (
	"fmt"
	"math/rand"
	"time"
)

// 用Select监听多个Channel
func ListenMultiWay() {
	ch1 := make(chan int, 1000)
	ch2 := make(chan byte, 1000)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch1 <- rand.Int()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch2 <- byte(rand.Int())
		}
	}()

AB:
	for {
		time.Sleep(time.Second)
		select { //同时监听多个channel，谁先有数据先执行哪个case，如果同时有数据则随机选一个case执行
		case v1 := <-ch1:
			fmt.Printf("v1=%d\n", v1)
		case v2 := <-ch2:
			fmt.Printf("v2=%d\n", v2)
			if v2 < 40 {
				break AB
			}
		default:
			fmt.Println("default")
		}
	}

	// 最后兜底，万一在break AB的瞬间，ch1里写入了新元素，则通过下面的select的把新元素打印出来。为什么要放在select里？因为如果ch1里没有元素我也不想阻塞，main协程想立即终止
	select { //无阻塞读channel的方式
	case v1 := <-ch1:
		fmt.Printf("at last v1=%d\n", v1)
	default:
	}
}

func SelectBlock() {
	go func() { // 工作协程
		for {
			time.Sleep(time.Second)
			fmt.Println("我还在对外提供服务")
		}
	}()

	select {} // 永久阻塞main协程，避免工作协程被强行终止
	// ch := make(chan struct{})
	// ch <- struct{}{}  // 永久阻塞
}
