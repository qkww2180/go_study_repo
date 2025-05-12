package main

import (
	"dqq/go/basic/type_func/oip/common"
	"dqq/go/basic/type_func/oip/filter"
	"dqq/go/basic/type_func/oip/recall"
	"dqq/go/basic/type_func/oip/sort"
	"fmt"
	"time"

	"log"
)

// 电商推荐
type Recommender struct {
	Recallers []recall.Recaller
	Sorter    sort.Sorter
	Filters   []filter.Filter
}

// 推荐主框架
func (rec *Recommender) Rec() []*common.Product {
	RecallMap := make(map[int]*common.Product, 100)
	//顺序执行多路召回
	for _, recaller := range rec.Recallers {
		begin := time.Now()
		products := recaller.Recall(10) //统一设置每路最多召回10个商品
		for _, product := range products {
			RecallMap[product.Id] = product //把多路召回的结果放到一个map里，按Id进行排重
		}
		log.Printf("召回%s用时%dns，召回%d个商品\n", recaller.Name(), time.Since(begin).Nanoseconds(), len(products))

	}
	log.Printf("排重后总共召回%d个商品\n", len(RecallMap))
	//把map转成slice
	RecallSlice := make([]*common.Product, 0, len(RecallMap))
	for _, product := range RecallMap {
		RecallSlice = append(RecallSlice, product)
	}
	//对召回的结果进行排序
	begin := time.Now()
	SortedResult := rec.Sorter.Sort(RecallSlice)
	log.Printf("排序%s用时%dns\n", rec.Sorter.Name(), time.Since(begin).Nanoseconds())
	//顺序执行多种过滤规则
	FilteredResult := SortedResult
	prevCount := len(FilteredResult)
	for _, filter := range rec.Filters {
		begin := time.Now()
		FilteredResult = filter.Filter(FilteredResult)
		log.Printf("过滤器%s用时%dns，过滤了%d个商品\n", filter.Name(), time.Since(begin).Nanoseconds(), prevCount-len(FilteredResult))
		prevCount = len(FilteredResult)
	}
	return FilteredResult
}

func main01() {
	rec := Recommender{
		//每种具体的实现可能是由不同的开发者完成。每种实现单独放一个文件，大家的代码互不干扰
		Recallers: []recall.Recaller{
			recall.HotRecall{Tag: "hot"},
			recall.SizeRecall{Tag: "size"},
		},
		Sorter: sort.RatioSorter{Tag: "ratio"},
		Filters: []filter.Filter{
			filter.AddressFilter{Tag: "address", City: "郑州"},
			filter.RatioFilter{Tag: "ratio"},
		},
	}
	// rec.Sorter = sort.SizeSorter{Tag: "size"}
	result := rec.Rec()
	for i, product := range result {
		fmt.Printf("第%d名：%d %s\n", i, product.Id, product.Name)
	}
}

//go run ./basic/oip
