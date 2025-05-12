package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

var (
	link = "https://www.bilibili.com/video/BV19T411x7Nj"
)

// 用标准库进行http GET请求
func stdGet() {
	if resp, err := http.Get(link); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			if bs, err := io.ReadAll(resp.Body); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("响应体")
				fmt.Println(string(bs))
			}
		}
	}
}

// HEAD类似于GET，返回的响应中没有具体的内容，用于获取报头。可用于快速验证一个URL是否可用
func checkHead() bool {
	if resp, err := http.Head(link); err != nil {
		fmt.Println(err) //如果link不存活，这里会报错
		return false
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 { //link虽然存活，但是请求不顺利
			return false
		}
		if len(resp.Header) > 0 {
			for k, v := range resp.Header {
				fmt.Println(k, v)
			}
			if bs, err := io.ReadAll(resp.Body); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("响应体")
				fmt.Println(string(bs))
			}
		} else {
			return false
		}
		return true
	}
}

func collyGet1() {
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("响应体")
		// fmt.Println(string(r.Body))
	})
	if err := c.Head(link); err == nil {
		if err := c.Visit(link); err != nil { //Visit执行的是GET请求。colly同样支持Post(URL string, requestData map[string]string)和PostRaw(URL string, requestData []byte)
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func collyGet2() {
	c := colly.NewCollector(colly.CheckHead())
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("响应体")
		fmt.Println(string(r.Body))
	})
	if err := c.Visit(link); err != nil { //Visit执行的是GET请求
		fmt.Println(err)
	}
}

func main1() {
	// if checkHead() {
	stdGet()
	// }

	// collyGet1()
	// collyGet2()
}
