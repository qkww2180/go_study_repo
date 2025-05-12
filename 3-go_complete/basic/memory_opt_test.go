package main

import (
	"testing"
)

func BenchmarkSynPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := Search("")
		_ = result //使用result
		pool.Put(result)
	}
}

// 16542             73337 ns/op    357630 B/op         19 allocs/op
// 75338             14542 ns/op     81921 B/op          1 allocs/op
// 284122              4724 ns/op         0 B/op          0 allocs/op
// 401432              2550 ns/op        24 B/op          1 allocs/op
