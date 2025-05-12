package main

import (
	"time"
)

func c1() {
	ch := make(chan int)
	ch <- 1
}

func c2() {
	ch := make(chan int)
	ch <- 1
	go func() {
		<-ch
	}()
}

func c3() {
	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
	}()
	ch <- 1
}

func c4() {
	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		<-ch
	}()
	ch <- 1
}

func main11() {
	c3()
}
