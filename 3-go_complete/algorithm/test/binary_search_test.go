package test

import (
	"dqq/algorithm"
	"math/rand"
	"sort"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	const L = 100 //数组的长度是L
LOOP:
	for i := 0; i < 100; i++ { //测试100个case
		arr := make([]int, L)
		for j := 0; j < L; j++ {
			arr[j] = rand.Intn(L)
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

		target := rand.Intn(L)
		index := algorithm.BinarySearch[int](arr, target)
		if index >= 0 {
			if arr[index] != target {
				t.Fail()
				break LOOP
			}
		} else {
			for _, ele := range arr {
				if ele == target {
					t.Fail()
					break LOOP
				}
			}
		}
	}
}

func TestBinarySearch1(t *testing.T) {
	arr := []int{1, 3, 6, 8, 10}
	targets := []int{1, 6, 10}
	for _, target := range targets {
		index := algorithm.BinarySearch[int](arr, target)
		if index > 0 {
			// fmt.Printf("case %d 测试失败\n", target)
			// t.Fail()
			// t.Errorf("case %d 测试失败\n", target)
			// break
			t.Fatalf("case %d 测试失败\n", target)
		}
	}
}

func TestBinarySearch2(t *testing.T) {

}

// go test ./algorithm/test -v -run=^TestBinarySearch$ -count=1
// go test ./algorithm/test -v -run=^TestBinarySearch\d+ -count=1
