package algorithm

import "cmp"

// 二分查找。返回target在arr中的下标，如果不存在则返回-1，如果target中存在多个只返回其中某一个的下标
func BinarySearch[T cmp.Ordered](arr []T, target T) int {
	begin := 0
	end := len(arr) - 1
	for begin <= end { //之所以是<=，而不是<，是因为区间内只剩下一个元素时也应该跟target进行比较，而不是直接返回-1
		middle := (begin + end) / 2
		if arr[middle] == target {
			return middle
		}
		if arr[middle] < target {
			begin = middle + 1
		} else {
			end = middle - 1
		}
	}
	return -1
}
