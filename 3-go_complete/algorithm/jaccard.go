package algorithm

import "cmp"

func JaccardTimeConsuming[T cmp.Ordered](collection1, collection2 []T) float64 {
	if len(collection1) == 0 || len(collection2) == 0 {
		return 0.0
	}
	map1 := make(map[T]struct{}, len(collection1))
	for _, ele := range collection1 {
		map1[ele] = struct{}{}
	}
	map2 := make(map[T]struct{}, len(collection2))
	for _, ele := range collection2 {
		map2[ele] = struct{}{}
	}
	intersection := 0 //交集的个数
	for key := range map1 {
		if _, exists := map2[key]; exists {
			intersection += 1
		}
	}
	return float64(intersection) / float64(len(collection1)+len(collection2)-intersection)
}

func JaccardForSorted[T cmp.Ordered](collection1, collection2 []T) float64 {
	if len(collection1) == 0 || len(collection2) == 0 {
		return 0.0
	}
	intersection := 0 //交集的个数
	for i, j := 0, 0; i < len(collection1) && j < len(collection2); {
		if collection1[i] == collection2[j] {
			intersection += 1
			i += 1
			j += 1
		} else if collection1[i] < collection2[j] {
			i += 1
		} else {
			j += 1
		}
	}
	return float64(intersection) / float64(len(collection1)+len(collection2)-intersection)
}
