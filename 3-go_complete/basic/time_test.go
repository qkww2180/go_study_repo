package main

import (
	"testing"
	"time"
)

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

// go test -v ./basic -bench=^Benchmark -count=1 -run=^$ -timeout=3s

/**
BenchmarkSleep-8             933           1293742 ns/op
BenchmarkAfter-8             943           1332043 ns/op
BenchmarkTimer-8             930           1330034 ns/op
BenchmarkTicker1-8          1198           1000106 ns/op
BenchmarkTicker2-8           943           1335797 ns/op
**/
