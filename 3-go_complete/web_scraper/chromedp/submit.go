// 登录cnblog
package main

import (
	"context"
	"dqq/web_scraper/util"
	"log"
	"os"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

func init() {
	config = util.CreateConfigReader("auth.yaml")
}

func main() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf), //是否开启debug日志
	)
	defer cancel()
	userInput := "#mat-input-0"
	passInput := "#mat-input-1"
	loginBtn := "body > app-root > app-sign-in-layout > div > div > app-sign-in > app-content-container > div > div > div > form > div > button"
	captcha := "#SM_BTN_WRAPPER_1"
	var content []byte
	login := chromedp.Tasks{
		chromedp.Navigate("https://account.cnblogs.com/signin"), //访问某页面
		chromedp.WaitVisible(loginBtn),
		// chromedp.Sleep(3 * time.Second),                               //等元素加载好
		chromedp.SetValue(userInput, config.GetString("cnblog.user")), //填入内容
		chromedp.SetValue(passInput, config.GetString("cnblog.pass")),
		chromedp.Click(loginBtn), //点击按钮
		chromedp.WaitVisible(captcha),
		chromedp.FullScreenshot(&content, 100), //截全屏,100是最高分辨率
	}
	if err := chromedp.Run(ctx, login); err != nil { //截屏数据保存到文件
		log.Fatal(err)
	}
	if err := os.WriteFile("z_data/after_login.png", content, 0o644); err != nil {
		log.Fatal(err)
	}
}

// go run ./web_scraper/chromedp
