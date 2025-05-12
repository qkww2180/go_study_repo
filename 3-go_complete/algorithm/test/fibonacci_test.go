package test

import (
	"dqq/algorithm"
	"fmt"
	"testing"
)

func TestFibonacci(t *testing.T) {
	n := 20
	a, b, c, d := algorithm.Fibonacci(n), algorithm.FibonacciTopDown(n), algorithm.FibonacciButtomUp(n), algorithm.FibonacciButtomUp_WithSpaceO1(n)
	if a != b || a != c || a != d {
		fmt.Println(a, b, c, d)
		t.Fail()
	}
}

func TestSteps(t *testing.T) {
	a := algorithm.Steps(10)
	b := algorithm.Steps(4)
	c := algorithm.Steps(5)
	if a != (c*c + b*b) {
		fmt.Println(a, b, c)
		t.Fail()
	}
}

func BenchmarkFibonacci(b *testing.B) {
	n := 20
	for i := 0; i < b.N; i++ {
		algorithm.Fibonacci(n)
	}
}

func BenchmarkFibonacciTopDown(b *testing.B) {
	n := 20
	for i := 0; i < b.N; i++ {
		algorithm.FibonacciTopDown(n)
	}
}

func BenchmarkFibonacciButtomUp(b *testing.B) {
	n := 20
	for i := 0; i < b.N; i++ {
		algorithm.FibonacciButtomUp(n)
	}
}

func BenchmarkFibonacciButtomUp_WithSpaceO1(b *testing.B) {
	n := 20
	for i := 0; i < b.N; i++ {
		algorithm.FibonacciButtomUp_WithSpaceO1(n)
	}
}

// go test ./algorithm/test -v -run=^TestFibonacci$ -count=1
// go test ./algorithm/test -v -run=^TestSteps$ -count=1
// go test ./algorithm/test -bench=^BenchmarkFibonacci -run=^$ -count=1 -benchmem -benchtime=2s

/**
BenchmarkFibonacci-8                               79546             26984 ns/op               0 B/op          0 allocs/op
BenchmarkFibonacciTopDown-8                     13425494               178.7 ns/op           176 B/op          1 allocs/op
BenchmarkFibonacciButtomUp-8                    42448996                58.46 ns/op          176 B/op          1 allocs/op
BenchmarkFibonacciButtomUp_WithSpaceO1-8        221522130               10.01 ns/op            0 B/op          0 allocs/op
*/
