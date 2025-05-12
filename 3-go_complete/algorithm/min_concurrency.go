package algorithm

import (
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
)

// 负载均衡算法--最小并发度
type MinimumConcurrency struct {
	endpoints   []string
	concurrency []int32 //每个endpoint对应的并发调用数
}

func NewMinimumConcurrency(endpoints []string) *MinimumConcurrency {
	concurrency := make([]int32, len(endpoints)) //len=cap=len(endpoints),初始切片里全是0
	return &MinimumConcurrency{
		endpoints:   endpoints,
		concurrency: concurrency,
	}
}

// 取出一个endpoint。返回的第一个值是endpoint的编号，归还时要使用这个编号
func (b *MinimumConcurrency) Take() (int, string) {
	//没有可用节点
	if len(b.endpoints) == 0 {
		return -1, ""
	}
	//只有一个节点，不需要做负载均衡
	if len(b.endpoints) == 1 {
		return 0, b.endpoints[0]
	}
	//选择concurrency最小的那个
	min := int32(math.MaxInt32)
	index := -1
	begin := rand.Intn(len(b.endpoints)) //使得endpoint被询问的顺序是随机的
	for i := 0; i < len(b.endpoints); i++ {
		idx := (i + begin) % len(b.endpoints)
		c := atomic.LoadInt32(&b.concurrency[idx])
		if c < min { //如果c=min，那当前的endpoint不会被选中。即concurrency相同的情况下，总是先询问的endpoint被选中，所以询问endpoint的顺序需要是随机的
			min = c
			index = idx
		}
	}
	atomic.AddInt32(&b.concurrency[index], 1) //发起请求时并发度加1
	return index, b.endpoints[index]
}

// 拿到请求结果后归还endpoint
func (b *MinimumConcurrency) Return(index int) error {
	if index >= len(b.endpoints) {
		return fmt.Errorf("index %d must less than %d", index, len(b.endpoints))
	}
	atomic.AddInt32(&b.concurrency[index], -1) //拿到请求结果后并发度减1
	return nil
}
