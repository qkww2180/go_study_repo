package main

import (
	"bytes"
	"dqq/web_scraper/util"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

func init() {
	config = util.CreateConfigReader("auth.yaml")
}

// 轮换socks5代理
func SetProxy() colly.CollectorOption {
	return func(c *colly.Collector) {
		proxies := make([]string, 0, 500)
		for _, ele := range util.ReadAllLines(util.RootPath + "/config/socks5.txt") {
			proxies = append(proxies, "socks5://"+ele)
		}
		rp, err := proxy.RoundRobinProxySwitcher(proxies...)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
}

// 并发&限速
func SetRateLimit(DomainGlob string, Parallelism int, RandomDelay int) colly.CollectorOption {
	return func(c *colly.Collector) {
		c.Limit(&colly.LimitRule{
			DomainGlob:  DomainGlob,
			Parallelism: Parallelism, //并发
			// Delay:       5 * time.Second,	//定长休息
			RandomDelay: time.Duration(RandomDelay) * time.Second, //随机休息
		})
	}
}

// 设置超时
func SetTimeout(DialTimeout, TLSHandshakeTimeout time.Duration) colly.CollectorOption {
	return func(c *colly.Collector) {
		c.WithTransport(&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   DialTimeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   TLSHandshakeTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		})
	}
}

// 访问失败时打印错误信息
func PrintError() colly.CollectorOption {
	return func(c *colly.Collector) {
		c.OnError(func(r *colly.Response, err error) {
			log.Println("Visiting", r.Request.URL, "failed:", err)
		})
	}
}

// 打印抓取到的完整html页面。Debug时使用，先确保能拿到html页面再去写解析html的代码
func PrintResponse() colly.CollectorOption {
	return func(c *colly.Collector) {
		c.OnResponse(func(r *colly.Response) {
			log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
		})
	}
}

// 随机切换User-Agent。colly自带的User-Agent数量是有限的
func RotateUserAgent() colly.CollectorOption {
	return func(c *colly.Collector) {
		extensions.RandomUserAgent(c) //生成一些逼真的Chrome/Firefox/Opera User-Agent
	}
}

type HttpHeader struct {
	Header map[string]string
}

type ScrapeopsResponse struct {
	Result []map[string]string `json:"result"`
}

var (
	headers = make(chan map[string]string, 30)
)

func init() {
	// SCRAPEOPS_KEY := config.GetString("scrapeops.key")
	// go func() {
	// 	for {
	// 		// 从scrapeops.io获取伪造的header，免费且无调用次数限制
	// 		if resp, err := http.Get("https://headers.scrapeops.io/v1/browser-headers?api_key=" + SCRAPEOPS_KEY + "&num_headers=10"); err != nil {
	// 			log.Printf("get header failed: %v", err)
	// 		} else {
	// 			defer resp.Body.Close()
	// 			if content, err := io.ReadAll(resp.Body); err == nil {
	// 				var sr ScrapeopsResponse
	// 				if err = sonic.Unmarshal(content, &sr); err != nil {
	// 					log.Printf("Unmarshal ScrapeopsResponse failed: %v", err)
	// 				} else {
	// 					for _, ele := range sr.Result {
	// 						headers <- ele
	// 					}
	// 				}
	// 			} else {
	// 				log.Printf("read scrapeops response failed: %v", err)
	// 			}
	// 		}
	// 	}
	// }()
}

// 随机生成https请求Header(包含User-Agent)
func RandomHeader() colly.CollectorOption {
	return func(c *colly.Collector) {
		c.OnRequest(func(r *colly.Request) {
			if header, ok := <-headers; ok {
				//注意：绝大部分浏览器在发送header时是按照固定的顺序的
				for k, v := range header {
					r.Headers.Set(k, v)
				}
			}
		})
	}
}

// 访问页面里的超链接
func VisitHref() colly.CollectorOption {
	return func(c *colly.Collector) {
		// 找到带有href属性的标签a
		c.OnHTML(`a[href]`, func(e *colly.HTMLElement) {
			link := e.Attr("href")
			c.Visit(e.Request.AbsoluteURL(link))
		})
	}
}
