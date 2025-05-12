package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	defautHeader = JwtHeader{
		Algo: "HS256",
		Type: "JWT",
	}
)

func login2(ctx *gin.Context) { //GIN框架默认为handler加了recover
	// userName := ctx.Query("user_name")
	// passWord := ctx.Query("pass_word")
	//根据userName查询数据库，确保passWord是正确的，并得到User的其他信息(比如用户昵称、用户角色、是否为vip等)
	header := defautHeader
	payload := JwtPayload{
		Issue:       "bilibili",
		IssueAt:     time.Now().Unix(),                         //因为每次的IssueAt不同，所以每次生成的token也不同
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(), //3天后过期，需要重新登录
		UserDefined: map[string]any{"name": "高性能golang", "role": "up主", "vip": true},
	}
	if token, err := GenJWT(header, payload); err != nil {
		log.Printf("生成token失败: %v", err)
		ctx.String(http.StatusInternalServerError, "token生成失败")
	} else {
		ctx.String(http.StatusOK, token)
	}
}

// JWT认证中间件
func jwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token") //从Header里取得token
		_, payload, err := VerifyJwt(token)
		if err != nil {
			ctx.String(http.StatusForbidden, "auth failed") //返回403
			ctx.Abort()                                     //通过Abort()使中间件后面的handler不再执行，但是本handler还是会继续执行。所以下一行代码需要显式return
			return
		}
		//token里包含了一些业务信息，可以放到context里往后传
		for k, v := range payload.UserDefined {
			ctx.Set(k, v)
		}
		ctx.Next()
	}
}

func main() {
	engine := gin.Default()
	engine.GET("/login", login2)
	engine.GET("/home", jwtAuthMiddleware(), myHomepage)
	engine.GET("/post", jwtAuthMiddleware(), postVideo)
	engine.Run("127.0.0.1:5678")
}

// curl -X GET '127.0.0.1:5678/login'
// curl -X GET '127.0.0.1:5678/home' --header 'token: xxxx'
// curl -X GET '127.0.0.1:5678/post' --header 'token: xxxx'
