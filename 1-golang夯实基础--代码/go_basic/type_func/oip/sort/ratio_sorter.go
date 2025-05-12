package sort

import (
	"dqq/go/basic/basic/oip/common"
	"sort"
)

type RatioSorter struct {
	Tag string
}

func (self RatioSorter) Name() string {
	return self.Tag
}

func (self RatioSorter) Sort(products []*common.Product) []*common.Product {
	sort.Slice(products, func(i, j int) bool {
		//按好评率降序排列
		return products[i].PositiveRatio > products[j].PositiveRatio
	})
	return products
}
