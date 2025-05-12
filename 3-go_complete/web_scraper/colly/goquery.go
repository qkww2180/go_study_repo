package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/sonic"
	"github.com/gocolly/colly/v2"
)

var (
	//传给goquery的字符需要是UTF-8编码
	html = `
	<html lang="en">
		<head>
			<meta charset="utf-8" />
			<title>员工信息</title>
		</head>
		<body>
			<div id="basic_div" class="gen">
				<h3>基本信息</h3>
				<ul>
					<li>姓名：<span class="bkl">大乔乔</span></li>
					<li>年龄：<span class="bkl">28</span></li>
				</ul>
				<div id="exp_div" class="gen">
				<h3>工作经历</h3>
				<p color="grey">高性能golang</p>
				</div>
			</div>
		</body>
	</html>
	`
)

func init() {
	blankReg := regexp.MustCompile(`[\n\t]`)
	html = blankReg.ReplaceAllString(html, "") //去除html中的所有换行符和制表符
}

// 解析html页面结构
func Document() {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html)) //接收一个io.Reader，也可以通过bufio.NewReader(*os.File)构建
	if err != nil {
		panic(err)
	}
	root := dom.Children().First()   //最外层的html标签
	language, _ := root.Attr("lang") //<html lang="en">
	fmt.Println(language)
	Level1Nodes := root.Children()
	fmt.Println(Level1Nodes.Length()) //<html>标签有几个孩子节点
	head := Level1Nodes.First()
	fmt.Println(head.Html())
	fmt.Println(head.Text())
	body := head.Siblings().First()
	div := body.Children().First()
	fmt.Println(div.HasClass("gen"), div.HasClass("red")) //<div id="basic_dic" class="gen">
	fmt.Println(strings.Repeat("=", 50))
}

// 定位html页面元素
func Selection() {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html)) //接收一个io.Reader，也可以通过bufio.NewReader(*os.File)构建
	if err != nil {
		panic(err)
	}
	allDiv := dom.Find("div")                       //根据标签名称选择元素
	fmt.Println(allDiv.Size())                      //页面中一共几个div
	allDiv.Each(func(i int, s *goquery.Selection) { //遍历这些div
		fmt.Println(s.Attr("id")) //打印这些div的id属性
	})
	fmt.Println()

	node := dom.Find("#basic_div") //根据id选择元素
	fmt.Println(node.Attr("id"))
	node = dom.Find(".gen")                       //根据class选择元素，会命中2个div
	node.Each(func(i int, s *goquery.Selection) { //遍历这些div
		fmt.Println(s.Attr("id")) //打印这些div的id属性
	})
	node = dom.Find("[color=grey]") //根据任意属性选择元素
	fmt.Println(node.Text())
	fmt.Println()

	node = dom.Find("p[color]") //寻找具有color属性的p标签
	fmt.Println(node.Text())
	node = dom.Find("p[color=grey]") //寻找color属性等于grey的p标签
	fmt.Println(node.Text())
	node = dom.Find("div[id][class=gen]") //寻找具有id属性且class=gen的div标签
	fmt.Println(node.Attr("id"))
	node = dom.Find("div#basic_div") //寻找id=basic_div的div标签
	fmt.Println(node.Attr("id"))
	fmt.Println()

	allDiv = dom.Find("body div")       //在body下寻找所有的div
	fmt.Println(allDiv.Size())          //找到2个
	div := dom.Find("body>div")         //在body的一级孩子节点中寻找div
	fmt.Println(div.Size())             //找到1个
	h := div.Find("ul li:nth-child(2)") //在div下找ul，再在ul下找li，li有多个，这里找第2个孩子节点（编号从1开始）
	fmt.Println(h.Text())               //年龄：28
	fmt.Println(strings.Repeat("=", 50))
}

type VideoDom struct {
	Title    string   `json:"title" selector:"h1.video-title"`
	Like     string   `json:"like" selector:"span.video-like-info.video-toolbar-item-text"`
	Coin     string   `json:"coin" selector:"span.video-coin-info.video-toolbar-item-text"`
	Favorite string   `json:"favorite" selector:"span.video-fav-info.video-toolbar-item-text"`
	Share    string   `json:"share" selector:"span.video-share-info.video-toolbar-item-text"`
	Keywords []string `json:"keywords" selector:"div.tag-panel a.tag-link"`
}

func colly1() {
	collector := colly.NewCollector()
	collector.OnHTML("div.left-container", func(h *colly.HTMLElement) {
		dom := h.DOM
		title := dom.Find("h1.video-title").Text()
		like := dom.Find("span.video-like-info.video-toolbar-item-text").Text()
		coin := dom.Find("span.video-coin-info.video-toolbar-item-text").Text()
		favorite := dom.Find("span.video-fav-info.video-toolbar-item-text").Text()
		share := dom.Find("span.video-share-info.video-toolbar-item-text").Text()
		keywords := make([]string, 0, 10)
		dom.Find("div.tag-panel a.tag-link").Each(func(i int, s *goquery.Selection) {
			keywords = append(keywords, s.Text())
		})

		video := VideoDom{
			Title:    title,
			Like:     like,
			Coin:     coin,
			Favorite: favorite,
			Share:    share,
			Keywords: keywords,
		}
		fmt.Println(sonic.MarshalString(video))
	})
	if err := collector.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
	fmt.Println(strings.Repeat("=", 50))
}

func colly2() {
	collector := colly.NewCollector()
	collector.OnHTML("div.left-container", func(h *colly.HTMLElement) {
		var video VideoDom
		if err := h.Unmarshal(&video); err == nil {
			fmt.Println(sonic.MarshalString(video))
		} else {
			fmt.Printf("unmarshal VideoDom failed: %v", err)
		}
	})
	if err := collector.Visit("https://www.bilibili.com/video/BV19T411x7Nj"); err != nil {
		fmt.Println(err)
	}
	fmt.Println(strings.Repeat("=", 50))
}

func main4() {
	Document()
	Selection()
	colly1()
	colly2()
}
