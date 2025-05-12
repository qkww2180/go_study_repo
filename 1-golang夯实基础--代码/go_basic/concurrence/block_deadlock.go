package concurrence

import (
	"fmt"
	"sync"
	"time"
)

func Block() {
	wg := sync.WaitGroup{}
	go func() {
		time.Sleep(5 * time.Second)
		// wg.Done()
	}()
	time.Sleep(time.Second) //时间到了，会自动解除阻塞
	fmt.Println("sleep over")

	wg.Add(1)
	wg.Wait()
	fmt.Println("wait over")

	ch := make(chan bool, 10)
	<-ch
	fmt.Println("receive channel over")
	for ele := range ch {
		fmt.Println(ele)
	}
	fmt.Println("traverse channel over")

	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
	fmt.Println("got lock")

	select {}

}
