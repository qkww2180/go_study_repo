package test

import (
	projectprepare "dqq/go/basic/b_project_prepare"
	"testing"
)

func TestCopySlice(t *testing.T) {

	var src, dest []int16
	src = []int16{1, 2, 3, 4}
	var c, n int

	c = len(src) - 1
	dest = make([]int16, c)
	n = projectprepare.CopySlice(dest, src)
	if n != c {
		t.Errorf("c=%d, n %d", c, n)
	}
	for i := 0; i < n; i++ {
		if dest[i] != src[i] {
			t.Errorf("c=%d, i=%d, dest %d src %d", c, i, dest[i], src[i])
		}
	}

	c = len(src)
	dest = make([]int16, c)
	n = projectprepare.CopySlice(dest, src)
	if n != len(src) {
		t.Errorf("c=%d, n %d", c, n)
	}
	for i := 0; i < n; i++ {
		if dest[i] != src[i] {
			t.Errorf("c=%d, i=%d, dest %d src %d", c, i, dest[i], src[i])
		}
	}

	c = len(src) + 1
	dest = make([]int16, c)
	n = projectprepare.CopySlice(dest, src)
	if n != len(src) {
		t.Errorf("c=%d, n %d", c, n)
	}
	// dest[0]--
	// dest[1]--
	for i := 0; i < n; i++ {
		if dest[i] != src[i] {
			t.Fatalf("c=%d, i=%d, dest %d src %d", c, i, dest[i], src[i])
		}
	}
}

func BenchmarkCopySlice(b *testing.B) {
	src := make([]int8, 10000)
	dest := make([]int8, 10000)
	b.ResetTimer() //开始计时
	for i := 0; i < b.N; i++ {
		projectprepare.CopySlice(dest, src)
	}
}

func BenchmarkStdCopySlice(b *testing.B) {
	src := make([]int8, 10000)
	dest := make([]int8, 10000)
	b.ResetTimer() //开始计时
	for i := 0; i < b.N; i++ {
		copy(dest, src)
	}
}

// go test -v ./b_project_prepare/test -run=TestCopySlice$ -count=1
// go test ./b_project_prepare/test -run=^$ -bench=CopySlice$ -count=1
