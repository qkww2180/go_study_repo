package main

import (
	"context"
	"fmt"
	"time"
)

var _ = func() {
L:
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if i%2 == 0 {
				break L
			}
		}
	}
}

var breakSwitch = func() {
LAB:
	for i := 0; i < 8; i++ {
		switch i {
		case 4:
			break LAB
		default:
			fmt.Println(i)
		}
	}
}

var breakSelect = func() {
	ch := make(chan int, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

L:
	for i := 0; i < 8; i++ {
		fmt.Println(i)
		select {
		case <-ch:
		case <-ctx.Done():
			break L
		}
	}
}

func main2() {
	breakSwitch()
	breakSelect()
}
