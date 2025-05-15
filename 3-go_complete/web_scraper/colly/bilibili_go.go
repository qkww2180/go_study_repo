package main

import (
	"bufio"
	"dqq/web_scraper/algorithm"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
	colly "github.com/gocolly/colly/v2"
)

var (
	videoBF     *algorithm.BloomFilter //把VideoID放入BloomFilter进行排重
	videoBFfile = "z_data/b_video.bf"
	bvidRegex   = regexp.MustCompile(`/video/(BV.+)/`)
	blankRegx   = regexp.MustCompile(`\s+`)
	goRegx      = regexp.MustCompile(`(^go\W+)|(\W+go$)|(\W+go\W+)`)
	golangRegx  = regexp.MustCompile(`(^golang\W+)|(\W+golang$)|(\W+golang\W+)`)

	fout    *os.File
	writer  *bufio.Writer
	err     error
	c       *colly.Collector
	success int

	spiderEntry     = "https://www.bilibili.com/video/BV1LS42197uW" //爬取的始点
	vedioCreateTime = "2024-01-01"                                  //只爬取vedioCreateTime之后发的视频
)

func init() {
	videoBF = algorithm.LoadBloomFilter(videoBFfile) //优先从文件中加载
	if videoBF == nil {
		videoBF = algorithm.NewBloomFilter(8, 80000) //如果文件中没有，则新创建
	}
}

func tearDown() {
	fout.Close()
	writer.Flush()
	if err := videoBF.Dump(videoBFfile); err != nil {
		log.Printf("dump url bloomfilter to file %s failed %v", videoBFfile, err)
	} else {
		log.Printf("dump url bloomfilter to file %s", videoBFfile)
	}
}

func main() {
	go listenSignal()

	c = colly.NewCollector(
		colly.AllowedDomains("www.bilibili.com"), //colly在抓取每个网址之前会先检查url是否满足正则、url是否已经抓取过
		colly.URLFilters(
			regexp.MustCompile(`.*/video/BV.+`),
		),
		RotateUserAgent(),
	)
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 2 * time.Second, //降低抓取频率，别让B站把你的IP封了
	})

	//抓取的内容存入CSV文件
	videoInfoFile := "z_data/b_go_video.csv"
	fout, err = os.OpenFile(videoInfoFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644) //append模式，程序重启后接着抓
	if err != nil {
		log.Printf("打开输出文件%s失败%s\n", videoInfoFile, err)
		return
	}
	writer = bufio.NewWriter(fout)
	defer tearDown()

	//从网页中抽取你需要的内容
	c.OnHTML("div.video-container-v1", scrapPage)
	if err := c.Visit(spiderEntry); err != nil {
		log.Printf("colly visit %s failed: %s", spiderEntry, err)
	}
}

// 字符串转数字
func toNumber(s string) int {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0
	}
	unit := 1.
	if strings.HasSuffix(s, "万") {
		unit = 10000.
		s = s[0 : len(s)-3]
	}
	if n, err := strconv.ParseFloat(s, 64); err != nil {
		log.Printf("parse number from %s failed: %s", s, err)
		return 0
	} else {
		return int(n * unit)
	}
}

// 从url抽取视频ID
func getBVID(url string) string {
	indexes := bvidRegex.FindAllStringSubmatchIndex(url, 1)
	if len(indexes) == 1 {
		if len(indexes[0]) == 4 {
			if indexes[0][2] >= 0 && indexes[0][3] > indexes[0][2] {
				return url[indexes[0][2]:indexes[0][3]]
			}
		}
	}
	return ""
}

// 接收到kill信号时把BloomFilter存入文件，再退出
func listenSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	tearDown()
	os.Exit(0)
}

// 从一个网页内抓取需要的内容
func scrapPage(h *colly.HTMLElement) {
	dom := h.DOM
	bvid := getBVID(h.Request.URL.Path)
	if len(bvid) == 0 {
		log.Printf("could not extract bvid from url %s", h.Request.URL)
		return
	}
	if videoBF.Exists(bvid) {
		log.Printf("vedio id %s has been exists", bvid)
		return
	}
	videoBF.Add(bvid) //抓取过的网页放入BloomFilter

	containGO := false
	title := strings.ToLower(dom.Find("h1.video-title").Text()) //转小写
	if goRegx.MatchString(title) || golangRegx.MatchString(title) {
		containGO = true
	}
	keywords := make([]string, 0, 10)
	dom.Find("div.tag-panel a.tag-link").Each(func(i int, s *goquery.Selection) {
		wd := blankRegx.ReplaceAllString(strings.ToLower(s.Text()), "")
		if goRegx.MatchString(wd) || golangRegx.MatchString(wd) {
			containGO = true
		}
		keywords = append(keywords, wd)
	})
	if !containGO { //视频标题和关键词中都不包含go，则跳过该网页
		// log.Printf("url %s not contain go", h.Request.URL)
		return
	}
	createTime := strings.TrimSpace(dom.Find("div.pubdate-ip-text").Text())
	if createTime < vedioCreateTime {
		// log.Printf("createTime %s", createTime)
		return
	}
	like := toNumber(dom.Find("span.video-like-info.video-toolbar-item-text").Text())
	coin := toNumber(dom.Find("span.video-coin-info.video-toolbar-item-text").Text())
	favorite := toNumber(dom.Find("span.video-fav-info.video-toolbar-item-text").Text())
	share := toNumber(dom.Find("span.video-share-info.video-toolbar-item-text").Text())
	view := toNumber(dom.Find("span.view.item").Text())
	upName := blankRegx.ReplaceAllString(dom.Find("a.up-name").Text(), "")
	writer.WriteString("https://www.bilibili.com/video/" + bvid)
	writer.WriteString("|")
	writer.WriteString(strings.ReplaceAll(title, "|", ""))
	writer.WriteString("|")
	writer.WriteString(createTime)
	writer.WriteString("|")
	writer.WriteString(upName)
	writer.WriteString("|")
	writer.WriteString(strconv.Itoa(view))
	writer.WriteString("|")
	writer.WriteString(strconv.Itoa(like))
	writer.WriteString("|")
	writer.WriteString(strconv.Itoa(coin))
	writer.WriteString("|")
	writer.WriteString(strconv.Itoa(favorite))
	writer.WriteString("|")
	writer.WriteString(strconv.Itoa(share))
	writer.WriteString("|")
	writer.WriteString(strings.Join(keywords, ","))
	writer.WriteString("\n")

	success++
	writer.Flush()
	log.Printf("deal video %s total %d", bvid, success) //打印进度

	recList := dom.Find("div.rec-list")
	hrefs := recList.Find("a[href]")
	hrefs.Each(func(i int, s *goquery.Selection) { //抓取推荐列表里的每一个视频
		link, _ := s.Attr("href")
		bvid := getBVID(link)
		if len(bvid) == 0 || videoBF.Exists(bvid) { //已经抓取过，则放弃本页面
			return
		}
		url := h.Request.AbsoluteURL(link)
		if err := c.Visit(url); err != nil {
			log.Printf("colly visit %s[src=https://www.bilibili.com%s] failed: %s", url, h.Request.URL.Path, err)
		}
	})
}

// go run .\web_scraper\colly\
