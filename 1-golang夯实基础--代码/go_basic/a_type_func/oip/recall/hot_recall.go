package recall

import (
	"dqq/go/basic/basic/oip/common"
	"sort"
)

var (
	allProducts = []*common.Product{
		{Id: 1, Name: "小勺", Size: 35, Sale: 3582, ShipAddress: "郑州", Price: 6.9, PositiveRatio: 0.76, RatioCount: 364},
		{Id: 2, Name: "袜子", Size: 23, Sale: 43654, ShipAddress: "郑州", Price: 6546, PositiveRatio: 0.86, RatioCount: 7},
		{Id: 3, Name: "裤子", Size: 354, Sale: 54, ShipAddress: "郑州", Price: 547, PositiveRatio: 0.96, RatioCount: 5436},
		{Id: 4, Name: "裙子", Size: 675, Sale: 756, ShipAddress: "郑州", Price: 5423, PositiveRatio: 0.86, RatioCount: 6452},
		{Id: 5, Name: "袖子", Size: 23, Sale: 423, ShipAddress: "郑州", Price: 64, PositiveRatio: 0.88, RatioCount: 235},
		{Id: 6, Name: "iPad", Size: 65, Sale: 3, ShipAddress: "郑州", Price: 87, PositiveRatio: 0.96, RatioCount: 254},
		{Id: 7, Name: "iPhone", Size: 146, Sale: 254, ShipAddress: "北京", Price: 143, PositiveRatio: 0.90, RatioCount: 254},
		{Id: 8, Name: "电脑", Size: 4, Sale: 543, ShipAddress: "北京", Price: 2354, PositiveRatio: 0.91, RatioCount: 5435},
		{Id: 9, Name: "机床", Size: 65, Sale: 8, ShipAddress: "北京", Price: 76, PositiveRatio: 0.44, RatioCount: 4213},
		{Id: 10, Name: "袜子", Size: 76, Sale: 143, ShipAddress: "北京", Price: 14, PositiveRatio: 0.68, RatioCount: 543},
		{Id: 11, Name: "裤子", Size: 67, Sale: 6354, ShipAddress: "北京", Price: 65, PositiveRatio: 0.89, RatioCount: 5436},
		{Id: 12, Name: "裙子", Size: 325, Sale: 4, ShipAddress: "北京", Price: 124, PositiveRatio: 0.85, RatioCount: 534},
	}
)

type HotRecall struct {
	Tag string
}

func (self HotRecall) Name() string {
	return self.Tag
}
func (self HotRecall) Recall(n int) []*common.Product {
	//按销量降序排序
	sort.Slice(allProducts, func(i, j int) bool {
		return allProducts[i].Sale > allProducts[j].Sale
	})
	//返回前n个
	rect := make([]*common.Product, 0, n)
	for i, ele := range allProducts {
		if i >= n {
			break
		}
		rect = append(rect, ele)
	}
	return rect
}
