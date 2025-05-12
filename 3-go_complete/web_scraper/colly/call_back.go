package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

type MyRequest struct {
	Url   string
	abort bool //取消网络请求
}

type MyResponse struct {
	Response    *http.Response
	ContentType string //返回什么类型的内容。html、xml、jpg、txt、zip等等
}

type Html struct {
}

func NewHtml(resp *MyResponse) *Html {
	return new(Html)
}

type Xml struct {
}

func NewXml(resp *MyResponse) *Xml {
	return new(Xml)
}

type MyCollector struct {
	ID                        uint32
	OnRequestCallBacks        []func(*MyRequest)
	OnResponseHeaderCallBacks []func(*MyResponse)
	OnErrorCallBacks          []func(*MyResponse, error)
	OnResponseCallBacks       []func(*MyResponse)
	OnHtmlCallBacks           []func(*Html)
	OnXmlCallBacks            []func(*Xml)
	OnScrapedCallBacks        []func(*MyResponse)
}

func (c *MyCollector) Visit(url string) error {
	request := &MyRequest{Url: url}
	for _, f := range c.OnRequestCallBacks {
		f(request)
	}
	if request.abort {
		return nil
	}

	resp, err := http.Get(request.Url)
	response := &MyResponse{Response: resp}
	if err != nil {
		for _, f := range c.OnErrorCallBacks {
			f(response, err)
		}
		return err
	} else {
		defer response.Response.Body.Close()
		for _, f := range c.OnResponseHeaderCallBacks {
			f(response)
		}
		if response.Response.StatusCode == 200 {
			for _, f := range c.OnResponseCallBacks {
				f(response)
			}

			switch response.ContentType {
			case "HTML":
				html := NewHtml(response)
				for _, f := range c.OnHtmlCallBacks {
					f(html)
				}
			case "XML":
				xml := NewXml(response)
				for _, f := range c.OnXmlCallBacks {
					f(xml)
				}
			}

			for _, f := range c.OnScrapedCallBacks {
				f(response)
			}
		} else {
			err = fmt.Errorf("response code %d", response.Response.StatusCode)
			for _, f := range c.OnErrorCallBacks {
				f(response, err)
			}
			return err
		}
	}
	return nil
}

func main2() {
	myCollector := new(MyCollector)
	myCollector.OnRequestCallBacks = []func(*MyRequest){
		func(r *MyRequest) {
			fmt.Printf("visit %s\n", r.Url)
			if strings.Index(r.Url, "baidu") >= 0 {
				r.abort = true
			}
		},
	}
	myCollector.OnResponseCallBacks = []func(*MyResponse){
		func(r *MyResponse) {
			fmt.Printf("got response from %s\n", r.Response.Request.URL)
		},
	}
	if err := myCollector.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
	if err := myCollector.Visit("https://www.baidu.com"); err != nil {
		fmt.Println(err)
	}
	fmt.Println()

	collector := colly.NewCollector()
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
