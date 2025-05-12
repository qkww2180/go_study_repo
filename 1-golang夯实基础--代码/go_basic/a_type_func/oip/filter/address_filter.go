package filter

import "dqq/go/basic/basic/oip/common"

type AddressFilter struct {
	Tag  string
	City string
}

func (self AddressFilter) Name() string {
	return self.Tag
}

func (self AddressFilter) Filter(products []*common.Product) []*common.Product {
	rect := make([]*common.Product, 0, len(products))
	for _, product := range products {
		if product.ShipAddress == self.City {
			rect = append(rect, product)
		}
	}
	return rect
}
