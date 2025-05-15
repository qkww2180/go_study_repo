package d_concurrence

import (
	"sync"
	"time"
)

var qps = make(chan struct{}, 100)

func handler() {
	qps <- struct{}{}
	defer func() {
		<-qps
	}()
	time.Sleep(3 * time.Second)
}

func QpsLimit() {
	const P = 1000
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			handler()
		}()
	}
	wg.Wait()
}
