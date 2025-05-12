package test

import (
	"bufio"
	"dqq/algorithm"
	"dqq/algorithm/util"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestTopKByPartition(t *testing.T) {
	const K = 10               //topK
	const L = 20               //数组长度
	for i := 0; i < 100; i++ { //测试100个case
		list := make([]int, L)
		for j := 0; j < L; j++ {
			list[j] = rand.Intn(5 * L)
		}
		topK := algorithm.TopKByPartition(list, K)

		min := topK[0] //topK里的最小值
		for _, ele := range topK {
			if min > ele {
				min = ele
			}
		}

		small, equal := 0, 0 //整个list里，比min小的和等于min的元素数
		for _, ele := range list {
			if ele < min {
				small++
			} else if ele == min {
				equal++
			}
		}

		//正常情况下，比min小的不可能超过len(list)-K个，比min小的加上等于min的至少应该等于len(list)-K+1
		if small > len(list)-K || small+equal < len(list)-K+1 {
			t.Fail()
		}
	}
}

func TestTopKByHeap(t *testing.T) {
	const K = 10
	const L = 20
	for i := 0; i < 100; i++ { //测试100个case
		list := make([]int, L)
		for j := 0; j < L; j++ {
			list[j] = rand.Intn(5 * L)
		}
		topK := algorithm.TopKByHeap(list, K)

		min := topK[0] //topK里的最小值
		for _, ele := range topK {
			if min > ele {
				min = ele
			}
		}

		small, equal := 0, 0 //整个list里，比min小的和等于min的元素数
		for _, ele := range list {
			if ele < min {
				small++
			} else if ele == min {
				equal++
			}
		}

		//正常情况下，比min小的不可能超过len(list)-K个，比min小的加上等于min的至少应该等于len(list)-K+1
		if small > len(list)-K || small+equal < len(list)-K+1 {
			t.Fail()
		}
	}
}

func BenchmarkTopKByHeap(b *testing.B) {
	const K = 10
	const L = 1000
	for i := 0; i < b.N; i++ {
		list := make([]int, L)
		for j := 0; j < L; j++ {
			list[j] = rand.Intn(5 * L)
		}
		algorithm.TopKByHeap(list, K)
	}
}

func BenchmarkTopKByPartition(b *testing.B) {
	const K = 10
	const L = 1000
	for i := 0; i < b.N; i++ {
		list := make([]int, L)
		for j := 0; j < L; j++ {
			list[j] = rand.Intn(5 * L)
		}
		algorithm.TopKByPartition(list, K)
	}
}

func TestFindFreqIpFromBigFile(t *testing.T) {
	topk1 := algorithm.FindFreqIpFromBigFile(util.RootPath+"data/ip_topk/ip.txt", 10)

	//用原始的方法统计出现次数最多的ip
	countMap := make(map[string]int, 10000)
	fin, err := os.Open(util.RootPath + "data/ip_topk/ip.txt")
	if err != nil {
		panic(err)
	}
	defer fin.Close()
	reader := bufio.NewReader(fin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					if v, exists := countMap[line]; exists {
						countMap[line] = v + 1
					} else {
						countMap[line] = 1
					}
				}
			}
			break
		} else {
			line = strings.TrimRight(line, "\n")
			if v, exists := countMap[line]; exists {
				countMap[line] = v + 1
			} else {
				countMap[line] = 1
			}
		}
	}
	list := make([]*algorithm.Item[int], 0, len(countMap))
	for k, v := range countMap {
		list = append(list, &algorithm.Item[int]{Info: k, Value: v})
	}
	sort.Slice(list, func(i, j int) bool { //通过排序找topK
		return list[j].Value < list[i].Value
	})

	//对比两种方法得到的结果
	for i, ele := range topk1 {
		fmt.Printf("%s:%d\t%s:%d\n", ele.Info, ele.Value, list[i].Info, list[i].Value)
	}
	fmt.Println()
}

// go test -v ./algorithm/test -run=^TestTopKByPartition$ -count=1
// go test -v ./algorithm/test -run=^TestTopKByHeap$ -count=1
// go test -v ./algorithm/test -run=^TestFindFreqIpFromBigFile$ -count=1
// go test ./algorithm/test -bench=^BenchmarkTopKByHeap$ -run=^$ -count=1 -benchmem -benchtime=5s
// go test ./algorithm/test -bench=^BenchmarkTopKByPartition$ -run=^$ -count=1 -benchmem -benchtime=5s

/**
BenchmarkTopKByHeap-8			343329		17251 ns/op		80 B/op		1 allocs/op
BenchmarkTopKByPartition-8		246532		20537 ns/op		8272 B/op	2 allocs/op
ByPartition和ByHeap时间复杂度相当，但ByPartition的空间复杂度是O(N)，ByHeap的空间复杂度是O(K)，所以ByPartition耗时更大一些
*/
