package filter

import "dqq/go/basic/basic/oip/common"

type RatioFilter struct {
	Tag string
}

func (self RatioFilter) Name() string {
	return self.Tag
}

func (self RatioFilter) Filter(products []*common.Product) []*common.Product {
	rect := make([]*common.Product, 0, len(products))
	for _, product := range products {
		if product.RatioCount > 10 && product.PositiveRatio > 0.8 {
			rect = append(rect, product)
		}
	}
	return rect
}
