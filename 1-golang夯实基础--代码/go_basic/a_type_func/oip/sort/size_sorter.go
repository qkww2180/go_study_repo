package sort

import (
	"dqq/go/basic/basic/oip/common"
	"sort"
)

type SizeSorter struct {
	Tag string
}

func (self SizeSorter) Name() string {
	return self.Tag
}

func (self SizeSorter) Sort(products []*common.Product) []*common.Product {
	sort.Slice(products, func(i, j int) bool {
		//按尺寸升序排列
		return products[i].Size > products[j].Size
	})
	return products
}
