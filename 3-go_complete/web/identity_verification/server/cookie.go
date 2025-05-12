package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

var (
	userInfos sync.Map //实践中使用Redis。这里用单机缓存代替，需要支持并发读写
)

const (
	authCookie = "auth"
)

// cookie name需要符合规则，否则该cookie会被Gin框架默默地丢弃掉
func genSessionId(ctx *gin.Context) string {
	return base64.StdEncoding.EncodeToString([]byte(ctx.Request.RemoteAddr))
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Vip  bool   `json:"vip"`
}

// 登录
func login(ctx *gin.Context) {
	// userName := ctx.Query("user_name")
	// passWord := ctx.Query("pass_word")
	//根据userName查询数据库，确保passWord是正确的，并得到User的其他信息(比如用户昵称、用户角色、是否为vip等)

	//为客户端生成cookie
	session_id := genSessionId(ctx)
	user := User{Name: "高性能golang", Role: "up主", Vip: true} //模拟从数据库中读取用户基本信息
	userInfo, _ := sonic.Marshal(user)
	//服务端维护所有客户端的cookie，用于对客户端进行认证，并缓存用户基本信息，减轻数据库的压力
	userInfos.Store(session_id, userInfo)
	//把cookie发给客户端。response header里会有一条 Set-Cookie: key=other，浏览器后续请求会自动把同域名下的cookie再放到request header里来，即request header里会有一条Cookie: key=other
	ctx.SetCookie( //cookie的大小不能超过4K
		authCookie,  //cookie name
		session_id,  //cookie value
		3000,        //maxAge，cookie的有效时间，时间单位秒。如果不设置过期时间，默认情况下关闭浏览器后cookie被删除
		"/",         //path，cookie存放目录
		"localhost", //cookie从属的域名,不区分协议和端口
		false,       //是否只能通过https访问。实际中敏感信息设为true
		false,       //是否允许别人通过js获取自己的cookie，设为false防止XSS攻击
	)
	fmt.Printf("set cookie %s = %s to client\n", authCookie, session_id)
	ctx.String(http.StatusOK, "登录成功")
}

// 访问我的主页
func myHomepage(ctx *gin.Context) {
	//some code here
	ctx.String(http.StatusOK, fmt.Sprintf("欢迎%s,这是你的个人主页", ctx.GetString("name"))) //从context里获得用户信息，避免读数据库
}

// 发布视频
func postVideo(ctx *gin.Context) {
	//some code here
	ctx.String(http.StatusOK, "视频发布成功")
}

func authMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session_id := genSessionId(ctx)
		//读取客户端的cookie
		sessionIdExists := false
		for _, cookie := range ctx.Request.Cookies() {
			if cookie.Name == authCookie && cookie.Value == session_id {
				sessionIdExists = true
				break
			}
		}
		//没传session_id，重定向到/go_to_login页面
		if !sessionIdExists {
			ctx.Redirect(http.StatusMovedPermanently, "http://localhost:5678/go_to_login") //重定向只能使用3xx的code
			// ctx.Abort()                                                                    //验证不通过，调用Abort，阻止下一个中间件(Handler)的执行
		}

		//验证Cookie Value是否正确
		if v, ok := userInfos.Load(session_id); !ok { //SessionId不存在
			fmt.Printf("session id %s 不存在\n", session_id)
			ctx.String(http.StatusForbidden, "身份认证失败")
			ctx.Abort() //验证不通过，调用Abort，阻止下一个中间件(Handler)的执行
		} else {
			var user User
			sonic.Unmarshal(v.([]byte), &user)
			ctx.Set("name", user.Name)
			ctx.Set("role", user.Role)
			ctx.Set("vip", user.Vip)
		}
	}
}

func main1() {
	engine := gin.Default()
	//搞一个静态html页面
	engine.LoadHTMLFiles("web/identity_verification/go_to_login.html")
	engine.GET("/go_to_login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "go_to_login.html", gin.H{})
	})

	engine.GET("/login", login)
	engine.GET("/home", authMiddleWare(), myHomepage)
	engine.GET("/post", authMiddleWare(), postVideo)
	engine.Run("localhost:5678")
}
