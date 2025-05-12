package recall

import (
	"dqq/go/basic/basic/oip/common"
)

type SizeRecall struct {
	Tag string
}

func (self SizeRecall) Name() string {
	return self.Tag
}
func (self SizeRecall) Recall(n int) []*common.Product {
	rect := make([]*common.Product, 0, n)
	for _, ele := range allProducts {
		if ele.Size < 200 { //只召回size小于200的商品
			rect = append(rect, ele)
			if len(rect) >= n {
				break
			}
		}
	}
	return rect
}
