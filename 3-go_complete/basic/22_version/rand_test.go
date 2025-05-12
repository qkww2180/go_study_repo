package main

import (
	"math/rand"
	rand2 "math/rand/v2"
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
		rand2.IntN(MAX)
	}
}

// go test ./basic/22_version -bench=Rand -run=^$ -count=1
