package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func timeMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()
		ctx.Next() //标准库是使用的next.ServeHTTP(rw, r)
		timeElapsed := time.Since(begin)
		msg := fmt.Sprintf("request %s use %d ms\n", ctx.Request.URL.Path, timeElapsed.Milliseconds())
		log.Printf(msg)
	}
}

func limitMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitCh <- struct{}{} //并发度达到100时就会阻塞
		log.Printf("concurrence %d\n", len(limitCh))
		ctx.Next()
		<-limitCh
	}
}

func boy(ctx *gin.Context) { //gin.Context里包含了http.ResponseWriter和*http.Request
	ctx.String(http.StatusOK, "hi boy ")
}

func main() {
	engine := gin.Default()
	engine.Use(timeMW()) //设置全局中间件。gin也支持按group设置中间件
	engine.Use(limitMW())
	engine.GET("/", boy)
	engine.Run("127.0.0.1:5678")
}
