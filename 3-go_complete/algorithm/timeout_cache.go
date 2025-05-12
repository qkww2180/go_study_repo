package algorithm

import (
	"container/heap"
	"time"
)

type HeapNode struct {
	value    int //对应到map里的key
	deadline int //到期时间戳，精确到秒
}

type TimeoutHeap []*HeapNode

func (heap TimeoutHeap) Len() int {
	return len(heap)
}
func (heap TimeoutHeap) Less(i, j int) bool {
	return heap[i].deadline < heap[j].deadline //小根堆
}
func (heap TimeoutHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}
func (heap *TimeoutHeap) Push(x interface{}) {
	node := x.(*HeapNode)
	*heap = append(*heap, node)
}
func (heap *TimeoutHeap) Pop() (x interface{}) {
	n := len(*heap)
	last := (*heap)[n-1]
	//删除最后一个元素
	*heap = (*heap)[0 : n-1]
	return last //返回最后一个元素
}

// 超时缓存。不支持并发
type TimeoutCache struct {
	cache map[int]interface{} //存储业务数据
	hp    TimeoutHeap         //辅助堆，用于决定淘汰谁
	cap   int                 //缓存容量的上限
}

func NewTimeoutCache(cap int) *TimeoutCache {
	tc := new(TimeoutCache)
	tc.cache = make(map[int]interface{}, cap)
	tc.hp = make(TimeoutHeap, 0, 10)
	tc.cap = cap
	heap.Init(&tc.hp) //包装升级，从一个常规的slice升级为堆
	go tc.taotai()
	return tc
}

// 向缓存中添加元素。时间复杂度 未到容量上限：O(1)  到达容量上限O(logN)
func (tc *TimeoutCache) Add(key int, value interface{}, life int) {
	//计算新元素的deadline
	deadline := int(time.Now().Unix()) + life
	if len(tc.cache) == tc.cap { //刚刚到达缓存容量上限
		top := tc.hp[0]
		if top.deadline <= deadline { //堆顶元素比新元素先到期。注意这里是<=，意味着在过期时间相等的情况下会保留新元素，淘汰已有的堆顶元素
			//淘汰堆顶
			heap.Pop(&tc.hp)
			delete(tc.cache, top.value)
		} else { //新元素比堆顶元素先到期，则新元素不能放入缓存
			return
		}
	}
	//直接把key value放入map
	tc.cache[key] = value
	//把key和deadline放入堆
	node := &HeapNode{value: key, deadline: deadline}
	heap.Push(&tc.hp, node)
}

// 从缓存中查找元素。时间复杂度O(1)
func (tc TimeoutCache) Get(key int) (interface{}, bool) {
	value, exists := tc.cache[key]
	return value, exists
}

// 不停地检查堆中是否有到期的元素，将其删除
func (tc *TimeoutCache) taotai() {
	//按deadline构建的小根堆，所以堆顶元素是最早到期的。每次for循环只需要检查堆顶就可以了，因为堆顶被删除后会立即产生新的堆顶
	for {
		if tc.hp.Len() == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		now := int(time.Now().Unix())
		top := tc.hp[0]
		if top.deadline < now { //堆顶的到期时间比当前时刻早，说明堆顶到期了
			heap.Pop(&tc.hp)
			delete(tc.cache, top.value)
		} else { //堆顶还没有到期
			time.Sleep(100 * time.Millisecond) //稍作休息，减少CPU开销
		}
	}
}
