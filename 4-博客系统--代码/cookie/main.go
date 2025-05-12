package main

import (
	"blog/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const COOKIE_NAME = "auth"

// var loggedIn = make(map[string]string, 1000) //维护已登录用户的信息

func Login(ctx *gin.Context) {
	//登录成功,给前端返回一个Set-Cookie
	uid := "8"
	key := util.RandStringRunes(20)
	// loggedIn[key] = uid
	SetCookieAuth(key, uid)
	fmt.Println("cookie value", key)

	//response header里会有一条 Set-Cookie: auth=xxx; other_key=other_value，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: auth=xxx; other_key=other_value
	ctx.SetCookie(
		COOKIE_NAME, //cookie的name
		key,         //cookie的value
		86400*7,     //cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
		"/",         //cookie存放目录
		"localhost", //cookie从属的域名,不区分协议和端口。如果不指定domain则默认为本host(如b.a.com)，如果指定的domain是一级域名(如a.com)，则二级域名(b.a.com)下也可以访问
		false,       //是否只能通过https访问
		true,        //是否允许别人通过js获取自己的cookie，设为false防止XSS攻击
	)
	ctx.String(http.StatusOK, "success")
}

func getUidFromCookie1(ctx *gin.Context) string {
	//http协议里没有cookie这个概念，cookie本质上是header里的一对KV
	// for _, cookie := range strings.Split(ctx.Request.Header.Get("cookie"), ";") {
	// 	arr := strings.Split(cookie, "=")
	// 	key := strings.TrimSpace(arr[0])
	// 	value := strings.TrimSpace(arr[1])
	// 	if key == COOKIE_NAME {
	// 		// if uid, exists := loggedIn[value]; exists {
	// 		if uid := GetCookieAuth(value); uid != "" {
	// 			return uid
	// 		}
	// 	}
	// }

	// go 1.23 新加了ParseCookie（针对request header中的Cookie）和ParseSetCookie（针对response header中的Set-Cookie）
	cookies, _ := http.ParseCookie(ctx.Request.Header.Get("Cookie"))
	for _, cookie := range cookies {
		key, value := cookie.Name, cookie.Value
		if key == COOKIE_NAME {
			// if uid, exists := loggedIn[value]; exists {
			if uid := GetCookieAuth(value); uid != "" {
				return uid
			}
		}
	}

	return ""
}
func getUidFromCookie2(ctx *gin.Context) string {
	//直接使用封装好的Request.Cookies()
	for _, cookie := range ctx.Request.Cookies() {
		fmt.Println(cookie.Name, cookie.Value)
		if cookie.Name == COOKIE_NAME {
			// if uid, exists := loggedIn[cookie.Value]; exists {
			if uid := GetCookieAuth(cookie.Value); uid != "" {
				return uid
			}
		}
	}
	return ""
}

func Account(ctx *gin.Context) {
	// if uid := getUidFromCookie1(ctx); uid != "" {
	ctx.HTML(http.StatusOK, "account.html", gin.H{"uid": ctx.GetString("uid")})
	// 	return
	// }
	// ctx.String(http.StatusForbidden, "请登录")
}

func Manage(ctx *gin.Context) {
	// if uid := getUidFromCookie1(ctx); uid != "" {
	ctx.HTML(http.StatusOK, "manage.html", gin.H{"uid": ctx.GetString("uid")})
	// 	return
	// }
	// ctx.String(http.StatusForbidden, "请登录")
}

func main() {
	engine := gin.Default()
	engine.LoadHTMLFiles("account.html", "manage.html")
	engine.GET("/login", Login)
	engine.GET("/account", Auth, Account)
	engine.GET("/manage", Auth, Manage)
	engine.Run("localhost:5678")
}

// cd cookie
// go run .
