package main

import "net/http"

/*
第一种方式：传统的面向接口编程
*/

// 搜索引擎类
type SearchEngine struct {
	Recallers []Recaller //召回，得到搜索结果
	Sorter    Sorter     //对召回结果进行排序
}

type Recaller interface {
	Recall() []int
}

type Sorter interface {
	Sort([]int) []int
}

// 具体的召回策略
func r() []int {
	return nil
}

// 具体的排序策略
func s([]int) []int {
	return nil
}

// Rec类实现了Recaller接口
type Rec struct{}

func (Rec) Recall() []int {
	return r()
}

func main2() {
	se := SearchEngine{
		Recallers: []Recaller{Rec{}},
	}
	_ = se

	se2 := SearchEngine2{
		Recallers: []func() []int{r},
		Sorter:    s,
	}
	_ = se2

	se = SearchEngine{
		Recallers: []Recaller{RecallType(r)},
		Sorter:    SortType(s),
	}
	_ = se

	//用http.Handle实现路由
	http.Handle("/", http.HandlerFunc(Boy)) //h_http.HandlerFunc类似于FT，它实现了http.Handler接口
}

/*
第2种方式：用函数替代接口
*/
type SearchEngine2 struct {
	Recallers []func() []int    //召回，得到搜索结果
	Sorter    func([]int) []int //对召回结果进行排序
}

/*
把函数转为接口的实现
*/
type RecallType func() []int
type SortType func([]int) []int

// RecallType类实现了Recaller接口
func (rt RecallType) Recall() []int {
	return rt()
}

// SortType类实现了Sorter接口
func (st SortType) Sort(arg []int) []int {
	return st(arg)
}

func Boy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("boy"))
}
