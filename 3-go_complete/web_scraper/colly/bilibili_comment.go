package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gocolly/colly/v2"
)

type CommentResponse struct {
	Code int `json:"code"`
	Data struct {
		Page struct {
			Pn    int `json:"num"`   //第几页
			Total int `json:"total"` //总共多少评论
		} `json:"page"`
		List []Comment `json:"list"`
	} `json:"z_data"`
}

type Comment struct {
	CreateTime int `json:"ctime"`
	Member     struct {
		UserId   string `json:"mid"`
		UserName string `json:"uname"`
		Gender   string `json:"sex"`
	} `json:"member"`
	Content struct {
		Message string `json:"message"`
	} `json:"content"`
	VedioTitle string `json:"title"`
	VedioID    string `json:"bvid"`
}

func main6() {
	c := colly.NewCollector()
	var total int //总共多少评论
	const pageSize = 10
	page := 1
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", config.GetString("bili.cookie")) //cookie里存着登录后的凭证
	})
	commentFile := "z_data/my_comment.csv"
	fout, err := os.OpenFile(commentFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Printf("打开输出文件%s失败%s\n", commentFile, err)
		return
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	c.OnResponse(func(r *colly.Response) {
		var resp CommentResponse
		if err := sonic.Unmarshal(r.Body, &resp); err == nil {
			total = resp.Data.Page.Total
			fmt.Println("total", total, "page", page)
			for _, ele := range resp.Data.List {
				if ele.Member.UserId == config.GetString("bili.myid") { //把我自己发的评论排除掉
					continue
				}
				ct := time.Unix(int64(ele.CreateTime), 0)
				writer.WriteString(ct.Format("2006-01-02 15:04:05"))
				writer.WriteString("|")
				writer.WriteString(ele.Member.UserId)
				writer.WriteString("|")
				writer.WriteString(ele.Member.UserName)
				writer.WriteString("|")
				writer.WriteString(ele.Member.Gender)
				writer.WriteString("|")
				writer.WriteString(strings.Trim(strings.Trim(ele.Content.Message, "\n"), "|"))
				writer.WriteString("|")
				writer.WriteString(ele.VedioTitle)
				writer.WriteString("|")
				writer.WriteString("https://www.bilibili.com/video/" + ele.VedioID)
				writer.WriteString("\n")
			}
		} else {
			log.Printf("Unmarshal CommentResponse failed: %s", err)
		}
	})

	for {
		url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/up/fulllist?order=1&filter=-1&type=1&bvid=&pn=%d&ps=%d&charge_plus_filter=false", page, pageSize)
		if err := c.Visit(url); err != nil {
			log.Printf("visit page %d failed:%s", page, err)
			break
		}
		if page*pageSize >= total {
			break
		}
		time.Sleep(5 * time.Second)
		page++
	}
}
