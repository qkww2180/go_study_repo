package d_concurrence

import (
	"fmt"
	"time"
)

// 通知其他所有人
func Broadcast() {
	ch := make(chan struct{})

	const P = 3
	for i := 0; i < P; i++ {
		go func() {
			<-ch
			fmt.Printf("%d 出发了\n", i)
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("大伙可以出发了")
	close(ch) //广播

	time.Sleep(time.Second)
}

// 等其他人都完成后，我再执行
func CountDownLatch() {
	const P = 3
	ch := make(chan struct{}, P)
	for i := 0; i < P; i++ {
		go func() {
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Printf("%d 完成工作了\n", i)
			ch <- struct{}{}
		}()
	}

	// 同步点。channel可以完成WaitGroup的功能
	for i := 0; i < P; i++ {
		<-ch
	}

	fmt.Println("其他人都执行完毕，我要开始了")
}
