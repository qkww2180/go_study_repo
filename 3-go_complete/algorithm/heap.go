package algorithm

import (
	"cmp"
	"errors"
	"slices"
)

// 自行实现一个堆
//
// go标准库中的堆 container/heap
type Heap[T cmp.Ordered] struct {
	arr []T
}

func NewHeap[T cmp.Ordered](arr []T) *Heap[T] {
	brr := slices.Clone(arr) //拷贝一份arr，别对arr造成影响
	return &Heap[T]{arr: brr}
}

// 向下调整
func (heap *Heap[T]) downwardAdjust(parent int) {
	left := 2*parent + 1       //左孩子在数组里的下标
	if left >= len(heap.arr) { //没有左孩子，说明parent已经是树的叶节点了
		return
	}

	//从父子三人（如果没有右孩子，则是父子二人）中找到最小者
	minIndex := parent
	minValue := heap.arr[minIndex]
	if heap.arr[left] < minValue {
		minValue = heap.arr[left]
		minIndex = left
	}
	right := 2*parent + 2 //右孩子在数组里的下标
	if right < len(heap.arr) {
		if heap.arr[right] < minValue {
			minIndex = right
		}
	}

	// 最小者跟父节点交换
	if minIndex != parent {
		heap.arr[minIndex], heap.arr[parent] = heap.arr[parent], heap.arr[minIndex]
		heap.downwardAdjust(minIndex) //递归向下调整。每当有元素调整下来时，要对以它为父节点的三角形区域进行调整
	}
}

// 构建堆
func (heap *Heap[T]) Build() {
	n := len(heap.arr)
	if n <= 1 {
		return
	}
	lastIndex := n / 2 * 2              //最下层、最右那个三角形区域的右孩子在数组里的下标
	for i := lastIndex; i > 0; i -= 2 { //逆序检查每一个三角形区域
		right := i
		parent := (right - 1) / 2   //根据右孩子的下标获得它的父节点的下标。右孩子可能是不存在的，但父节点一定存在
		heap.downwardAdjust(parent) //调整该三角形区域
	}
}

// 向上调整
func (heap *Heap[T]) upwardAdjust(currIndex int) {
	if currIndex == 0 { //已经来到最顶层
		return
	}
	parent := (currIndex - 1) / 2
	// 如果currIndex比父节点小，则两者交换，然后递归向下；否则函数退出
	if heap.arr[currIndex] < heap.arr[parent] {
		heap.arr[currIndex], heap.arr[parent] = heap.arr[parent], heap.arr[currIndex]
		heap.upwardAdjust(parent)
	}
}

// 插入一个元素
func (heap *Heap[T]) Push(ele T) {
	heap.arr = append(heap.arr, ele)     //直接向数组尾部追加一个元素
	heap.upwardAdjust(len(heap.arr) - 1) //向上调整
}

// 删除并返回堆顶元素
func (heap *Heap[T]) Pop() (T, error) {
	if len(heap.arr) == 0 {
		var v T
		return v, errors.New("heap is empty")
	}
	root := heap.arr[0]
	//直接用最后一个元素替换首元素
	heap.arr[0] = heap.arr[len(heap.arr)-1]
	heap.arr = heap.arr[:len(heap.arr)-1] //删除最后一个元素
	//向下调整
	heap.downwardAdjust(0)
	return root, nil
}

// 取得堆顶元素
func (heap *Heap[T]) Top() (T, error) {
	if len(heap.arr) == 0 {
		var v T
		return v, errors.New("heap is empty")
	}
	return heap.arr[0], nil
}

// 堆里面有几个元素
func (heap *Heap[T]) Size() int {
	return len(heap.arr)
}

// 获取堆里的所有元素
func (heap *Heap[T]) GetAll() []T {
	return heap.arr
}

// 替换堆顶元素
func (heap *Heap[T]) ReplaceTop(ele T) {
	if len(heap.arr) == 0 {
		heap.Push(ele)
	} else {
		heap.arr[0] = ele
		//向下调整
		heap.downwardAdjust(0)
	}
}
