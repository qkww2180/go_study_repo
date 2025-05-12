package algorithm

import (
	"container/list"
	"log"
	rand "math/rand/v2"
)

type AliasSampler struct {
	accept []float64 //属于自己原本的概率是多少
	alias  []int     //上半部分是从谁那儿填补过来的
}

// 构建采样器。核心是构建数组accept和alias
func NewAliasSampler(probs []float64) *AliasSampler {
	if len(probs) <= 0 {
		log.Printf("invalid arguments")
		return nil
	}
	//先将概率归一化
	sum := 0.0
	for _, ele := range probs {
		sum += ele
	}
	for i, ele := range probs {
		probs[i] = ele / sum
	}
	// 每个事件的概率乘以n
	n := float64(len(probs))
	arr := make([]float64, len(probs))
	for i, ele := range probs {
		arr[i] = n * ele
	}

	accept := make([]float64, len(probs))
	alias := make([]int, len(probs))
	//初始化accept和alias，默认没有从别人那儿借概率
	for i := 0; i < len(probs); i++ {
		accept[i] = 1.0
		alias[i] = -1
	}
	largeStack := list.New()
	smallStack := list.New()
	//概率小于1的放入small栈，大于1的放入large栈
	for i, ele := range arr {
		if ele < 1 {
			smallStack.PushFront(i)
		} else {
			largeStack.PushFront(i)
		}
	}
	// 将面积大于1的事件多出的面积补充到面积小于1对应的事件中，以确保每一个小方格的面积为1
	for largeStack.Len() > 0 && smallStack.Len() > 0 {
		//分别从2个栈顶取出一个元素
		smallIndex := smallStack.Front()
		smallStack.Remove(smallIndex)
		largeIndex := largeStack.Front()
		largeStack.Remove(largeIndex)
		largeIdx := largeIndex.Value.(int)
		smallIdx := smallIndex.Value.(int)
		//属于自己原本的概率是多少
		accept[smallIdx] = arr[smallIdx]
		//上半部分是从谁那儿填补过来的
		alias[smallIdx] = largeIdx
		//large的补给small了，large的概率要减小
		arr[largeIdx] -= (1 - arr[smallIdx])
		//smallIdx已经补够1了，不再small了；但largeIdx还得放回栈中
		if arr[largeIdx] < 1.0 {
			smallStack.PushFront(largeIdx)
		} else if arr[largeIdx] > 1.0 {
			largeStack.PushFront(largeIdx)
		}
	}
	return &AliasSampler{
		accept: accept,
		alias:  alias,
	}
}

// 采样。支持并发调用
func (as AliasSampler) Sample() int {
	i := rand.IntN(len(as.alias)) //生成[0,N-1]上的随机整数i
	f := rand.Float64()           //生成U(0,1)上的随机小数f
	if f < as.accept[i] {
		return i //如果f<accept[i]则采样i
	} else {
		return as.alias[i] //否则采样alias[i]
	}
}
