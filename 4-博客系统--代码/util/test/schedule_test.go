package test

import (
	"blog/util"
	"log"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	log.Println("begin")
	time.Sleep(1 * time.Second)
	log.Println("sleep")
	<-time.After(1 * time.Second)
	log.Println("time.After")
	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()
	<-timer.C
	log.Println("timer")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for i := 0; i < 3; i++ {
		<-ticker.C
		// ticker.Reset(1 * time.Second)
		log.Println("ticker")
	}
}

func BenchmarkSleep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Sleep(1 * time.Millisecond)
	}
}

func BenchmarkAfter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		<-time.After(1 * time.Millisecond)
	}
}

func BenchmarkTimer(b *testing.B) {
	timer := time.NewTimer(1 * time.Millisecond)
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		<-timer.C
		timer.Reset(1 * time.Millisecond) //timer只能使用一次，想再次使用必须Reset
	}
}

func BenchmarkTicker1(b *testing.B) {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()
	for i := 0; i < b.N; i++ {
		<-ticker.C //Ticker可以一直使用
	}
}

func BenchmarkTicker2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ticker := time.NewTicker(1 * time.Millisecond)
		<-ticker.C
		ticker.Stop()
	}
}

func TestSchedule(t *testing.T) {
	go util.Shedule(func() {
		log.Println("schedule 1")
		time.Sleep(70 * time.Second)
	}, 0, 0, 0, true)
	go util.Shedule(func() {
		log.Println("schedule 2")
		time.Sleep(70 * time.Second)

	}, 10, 4, 20, false)
	select {} //永远阻塞
}

// go test -v ./util/test -run=^TestWait$ -count=1
// go test -v ./util/test -run=^TestSchedule$ -count=1 -timeout=100m   //timeout设大一点
// go test -v ./util/test -bench=^Benchmark -count=1 -run=^$ -timeout=3s

/**
BenchmarkSleep-8             933           1293742 ns/op
BenchmarkAfter-8             943           1332043 ns/op
BenchmarkTimer-8             930           1330034 ns/op
BenchmarkTicker1-8          1198           1000106 ns/op
BenchmarkTicker2-8           943           1335797 ns/op
**/
