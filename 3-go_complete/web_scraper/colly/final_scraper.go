package main

import (
	"fmt"
	"log"
	"regexp"
	"sync/atomic"

	"github.com/bytedance/sonic"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/redisstorage"
)

var cnt int32

func write2DB(video *VideoDom) {
	fmt.Println(sonic.MarshalString(video))
	if atomic.AddInt32(&cnt, 1)%10 == 0 {
		log.Printf("已经抓取了%d个网页", atomic.LoadInt32(&cnt))
	}
}

func main5() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.bilibili.com"), //colly在抓取每个网址之前会先检查url是否满足正则、url是否已经抓取过
		colly.URLFilters(
			regexp.MustCompile(`.*/video/BV.+`),
		),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 3, //并行抓取
	})
	//Cookie和访问过的URL默认都存储在本地内存中，也可以指定存到Redis、SQLite3或MongoDB中
	//如果本地没启redis，把storage代码注释掉
	storage := &redisstorage.Storage{
		Address:  "127.0.0.1:6379",
		Password: "",
		DB:       0,
		Prefix:   "bili_storage",
	}
	// 设置Storage
	if err := c.SetStorage(storage); err != nil {
		panic(err)
	}
	// 如果需要，可以先清空Storage
	if err := storage.Clear(); err != nil {
		panic(err)
	}
	defer storage.Client.Close()

	//发生错误时，打印错误信息
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Visiting", r.Request.URL, "failed:", err)
	})
	//有多个OnHTML回调函数时注意它们是按照注册的顺序依次执行的，所以根据href继续visit的OnHTML要放在最后
	c.OnHTML("div.left-container", func(h *colly.HTMLElement) {
		var video VideoDom
		if err := h.Unmarshal(&video); err == nil {
			write2DB(&video)
		} else {
			fmt.Printf("unmarshal VideoDom failed: %v", err)
		}
	})
	c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link)) //一般href都是路径，通过AbsoluteURL转成完整的url
	})

	if err := c.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
}

// go run ./web_scraper/colly
