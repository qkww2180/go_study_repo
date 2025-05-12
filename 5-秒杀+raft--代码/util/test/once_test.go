package test

import (
	"dqq/concurrency/util"
	"sync"
	"testing"
)

func TestOnce(t *testing.T) {
	var a int
	once := util.Once{}
	const P = 10
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			once.Do(func() {
				a++
			})
		}()
	}
	wg.Wait()
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			once.Do(func() {
				a++
			})
		}()
	}
	wg.Wait()
	if a != 1 {
		t.Fatalf("a=%d", a)
	}
}

func BenchmarkMyOnce(b *testing.B) {
	var a int
	once := util.Once{}
	for i := 0; i < b.N; i++ {
		once.Do(func() {
			a++
		})
	}
}

func BenchmarkStdOnce(b *testing.B) {
	var a int
	once := sync.Once{}
	for i := 0; i < b.N; i++ {
		once.Do(func() {
			a++
		})
	}
}

var g int
var m sync.Mutex

func fast() {
	if g == 0 {
		slow(&g)
	}
}

func slow(g *int) {
	m.Lock()
	defer m.Unlock()
	arr := make([]int, 1024)
	set := make(map[int]struct{}, 1024)
	for _, ele := range arr {
		set[ele] = struct{}{}
	}
	*g = 1
}

func inline() {
	if g == 0 {
		m.Lock()
		defer m.Unlock()
		arr := make([]int, 1024)
		set := make(map[int]struct{}, 1024)
		for _, ele := range arr {
			set[ele] = struct{}{}
		}
		g = 1
	}
}

func BenchmarkFastSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fast()
	}
}

func BenchmarkInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inline()
	}
}

func TestReadWrite(t *testing.T) {
	a := uint64(1<<63) - 1
	b := a
	const P = 10000
	wg := sync.WaitGroup{}
	wg.Add(2 * P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			if b != a {
				t.Errorf("b=%d", b)
			}
		}()
		go func() {
			defer wg.Done()
			b = a
		}()
	}
	wg.Wait()
}

// go test -v ./util/test -run=^TestOnce$ -count=1
// go test ./util/test -bench=Once$ -run=^$ -count=1 -benchmem
// go test ./util/test -bench=BenchmarkFastSlow$ -run=^$ -count=1 -benchmem
// go test ./util/test -bench=BenchmarkInline$ -run=^$ -count=1 -benchmem
// go test -v ./util/test -run=^TestReadWrite$ -count=1
