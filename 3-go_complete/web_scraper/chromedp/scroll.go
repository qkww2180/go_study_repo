// 下载百度图片
package main

import (
	"context"
	"dqq/web_scraper/util"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var (
	imgCh        = make(chan string, 1000)
	imgID        int32
	downloadTask = make(chan struct{})
)

func download() {
	//如果存储图片的目录不存在，就先创建该目录
	imgFolder := util.RootPath + "z_data/baidu_img"
	if exists, _ := util.PathExists(imgFolder); !exists {
		os.MkdirAll(imgFolder, 0o644)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	for {
		if urlstr, ok := <-imgCh; !ok {
			break
		} else {
			done := make(chan bool)
			var requestID network.RequestID
			chromedp.ListenTarget(ctx, func(v interface{}) {
				switch ev := v.(type) {
				case *network.EventRequestWillBeSent: //发送请求
					if ev.Request.URL == urlstr {
						requestID = ev.RequestID //记下requestID，确保该request请求的确实是urlstr
					}
				case *network.EventLoadingFinished:
					if ev.RequestID == requestID {
						close(done) //加载完毕，关闭channel
					}
				}
			})

			if err := chromedp.Run(ctx,
				chromedp.Navigate(urlstr), //访问图片
			); err != nil {
				log.Fatal(err)
			}

			<-done //等待加载完毕
			var buf []byte
			if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
				var err error
				buf, err = network.GetResponseBody(requestID).Do(ctx) //获取response body,放到buf里
				return err
			})); err != nil {
				log.Fatal(err)
			}

			//把buf保存到文件
			if err := os.WriteFile(imgFolder+strconv.Itoa(int(atomic.AddInt32(&imgID, 1)))+".jpg", buf, 0o644); err != nil {
				log.Fatal(err)
			}
		}
	}
	downloadTask <- struct{}{}
}

func main2() {
	go download()
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf), //是否开启debug日志
	)
	defer cancel()
	url := "https://image.baidu.com/search/albumsdetail?tn=albumsdetail&word=%E5%9F%8E%E5%B8%82%E5%BB%BA%E7%AD%91%E6%91%84%E5%BD%B1%E4%B8%93%E9%A2%98&fr=searchindex_album%20&album_tab=%E5%BB%BA%E7%AD%91&album_id=7&rn=30"
	var nodes []*cdp.Node
	chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("#loadingBar"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			//滚动10次
			for i := 0; i < 10; i++ {
				if err := chromedp.Run(ctx,
					// 滚动到页面底部
					chromedp.Evaluate(`window.scrollTo(0,document.body.scrollHeight);`, nil),
					// 暂停1秒钟
					chromedp.Sleep(1*time.Second),
				); err != nil {
					return err
				}
			}

			//向下箭头按1000次
			// for i := 0; i < 1000; i++ {
			// 	if err := chromedp.Run(ctx, chromedp.KeyEvent(kb.ArrowDown)); err != nil {
			// 		return err
			// 	}
			// }

			return nil
		}),
		chromedp.Nodes("a.albumsdetail-item", &nodes, chromedp.ByQueryAll), //定位到页面上所有的图片

	)

	for i, node := range nodes {
		var src string
		chromedp.Run(ctx, chromedp.AttributeValue("z_img", "src", &src, nil, chromedp.ByQuery, chromedp.FromNode(node))) //把图片URL放入imgCh
		imgCh <- src
		fmt.Println(i, src)
	}
	close(imgCh)
	<-downloadTask
}

// go run ./web_scraper/chromedp
