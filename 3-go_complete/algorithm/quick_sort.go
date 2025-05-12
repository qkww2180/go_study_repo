package algorithm

import "cmp"

//原地快速排序
func Partition[T cmp.Ordered](arr []T) {
	if len(arr) <= 1 { //递归终止条件
		return
	}
	pivot := 0
	i := 1
	j := len(arr) - 1
	for i < j { //一直循环，直到i和j相遇
		//先移动j
		for ; i < j; j-- {
			if arr[j] < arr[pivot] { // slice[0]作为基准
				break
			}
		}
		//后移动i
		for ; i < j; i++ {
			if arr[i] > arr[pivot] {
				break
			}
		}
		//i和j还没有相遇，则交换slice[i]和slice[j]
		if i < j {
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	//i和j相遇后，如果slice[i] < slice[pivot]，则i和pivot需要交换
	if arr[i] < arr[pivot] {
		//交换slice[pivot]和slice[i]
		arr[pivot], arr[i] = arr[i], arr[pivot]
		pivot = i
	}
	//在pivot左右两侧的两个子切片上递归调用partition
	if pivot > 0 {
		Partition(arr[:pivot])
	}
	if pivot < len(arr)-1 {
		Partition(arr[pivot+1:])
	}
}
