package main_test

import (
	"fmt"
	"math/rand"
	randv2 "math/rand/v2"
	"testing"
)

const MAX = 1e9

func BenchmarkRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(MAX)
	}
}

func BenchmarkRandV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randv2.IntN(MAX)
	}
}

func TestRand(t *testing.T) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d  ", randv2.IntN(100))
	}
}

func TestRandSeed(t *testing.T) {
	source := randv2.NewPCG(123, 456)
	for i := 0; i < 5; i++ {
		// source.Seed(123, 456)
		rander := randv2.New(source)
		fmt.Printf("%d  ", rander.IntN(100))
	}
}

// go test -v ./type_func -run=^TestRand$ -count=1
// go test -v ./type_func -run=^TestRandSeed$ -count=1
// go test ./type_func -bench=Rand -run=^$ -count=1
/*
BenchmarkRand-8         100000000               11.43 ns/op
BenchmarkRandV2-8       182933689                6.206 ns/op
*/
