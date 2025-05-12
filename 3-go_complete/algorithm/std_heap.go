package algorithm

import "cmp"

type Item[T cmp.Ordered] struct {
	Info  string
	Value T
}

type PriorityQueue[T cmp.Ordered] []*Item[T]

func (pq PriorityQueue[T]) Len() int {
	return len(pq)
}

func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].Value < pq[j].Value //golang默认提供的是小根堆，而优先队列是大根堆，所以这里要反着定义Less
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// 往slice里append,需要传slice指针
func (pq *PriorityQueue[T]) Push(x interface{}) {
	item := x.(*Item[T])
	*pq = append(*pq, item)
}

// 让slice指向新的子切片，需要传slice指针
func (pq *PriorityQueue[T]) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]   //数组最后一个元素
	*pq = (*pq)[0 : n-1] //去掉最一个元素
	return item
}
