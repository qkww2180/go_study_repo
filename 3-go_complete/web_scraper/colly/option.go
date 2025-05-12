package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func (c *MyCollector) OnRequest(f func(*MyRequest)) {
	if c.OnRequestCallBacks == nil {
		c.OnRequestCallBacks = make([]func(*MyRequest), 0, 5)
	}
	c.OnRequestCallBacks = append(c.OnRequestCallBacks, f)
}

type ResponseCallBack func(*MyResponse)

func (c *MyCollector) OnResponse(f ResponseCallBack) {
	if c.OnResponseCallBacks == nil {
		c.OnResponseCallBacks = make([]func(*MyResponse), 0, 5)
	}
	c.OnResponseCallBacks = append(c.OnResponseCallBacks, f)
}

// 一个函数以struct指针为参数，把这个函数定义为一种类型，该类型通常以Option结尾
type MyCollectorOption func(*MyCollector)

func ID(id uint32) MyCollectorOption {
	return func(c *MyCollector) {
		c.ID = id
	}
}

// 构建函数中以Option作为不定长参数
func NewMyCollector(options ...MyCollectorOption) *MyCollector {
	c := new(MyCollector)
	for _, opt := range options {
		opt(c)
	}
	return c
}

func main3() {
	myCollector := NewMyCollector(
		ID(4),
	)
	fmt.Printf("collector id %d\n", myCollector.ID)
	myCollector.OnRequest(func(r *MyRequest) {
		fmt.Printf("visit %s\n", r.Url)
		if strings.Index(r.Url, "baidu") >= 0 {
			r.abort = true
		}
	})
	myCollector.OnResponse(func(r *MyResponse) {
		fmt.Printf("got response from %s\n", r.Response.Request.URL)
	})
	if err := myCollector.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
	if err := myCollector.Visit("https://www.baidu.com"); err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	collector := colly.NewCollector(
		colly.ID(4),
	)
	fmt.Printf("collector id %d\n", collector.ID)
	collector.OnRequest(func(r *colly.Request) {
		fmt.Printf("visit %s\n", r.URL)
		if strings.Index(r.URL.Host, "baidu") >= 0 {
			r.Abort()
		}
	})
	collector.OnResponse(func(r *colly.Response) {
		fmt.Printf("got response from %s\n", r.Request.URL)
	})
	if err := collector.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
	if err := collector.Visit("https://www.baidu.com"); err != nil {
		fmt.Println(err)
	}
}
