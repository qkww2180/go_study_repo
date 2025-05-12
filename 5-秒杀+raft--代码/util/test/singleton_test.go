package test

import (
	"dqq/concurrency/util"
	"fmt"
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {
	cfg := util.GetConfig()
	fmt.Println(cfg.Password, cfg.ServerAddress)

	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			util.GetConfig()
		}()
	}
	wg.Wait()
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			util.GetConfig()
		}()
	}
	wg.Wait()
}

// go test -v ./util/test -run=^TestSingleton$ -count=1
