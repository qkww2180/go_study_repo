package algorithm

import (
	"cmp"
	"fmt"
)

type Comparable interface {
	int | byte | ~int32 | ~int64 | string | float32 | float64 // ~T: ~ 是Go 1.18新增的符号，~T表示底层类型是T的所有类型
}

// 找出众多元素的最小者
func Min[T cmp.Ordered](list ...T) (T, error) {
	var rect T
	if len(list) == 0 {
		return rect, fmt.Errorf("should not find min element from empty collectoin")
	}
	rect = list[0]
	for i := 1; i < len(list); i++ {
		if list[i] < rect {
			rect = list[i]
		}
	}
	return rect, nil
}

// 找出众多元素的最大者
func Max[T Comparable](list ...T) (T, error) {
	var rect T
	if len(list) == 0 {
		return rect, fmt.Errorf("should not find max element from empty collectoin")
	}
	rect = list[0]
	for i := 1; i < len(list); i++ {
		if list[i] > rect {
			rect = list[i]
		}
	}
	return rect, nil
}
