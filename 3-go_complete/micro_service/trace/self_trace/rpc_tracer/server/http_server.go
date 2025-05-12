package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func greet(ctx *gin.Context) {
	time.Sleep(100 * time.Millisecond)
	ctx.String(http.StatusOK, "hello")
}

func traceMiddleware(ctx *gin.Context) {
	begin := time.Now()

	//从request header里抽取trace信息
	var traceId string = ctx.GetHeader("trace_id")
	var userId string = ctx.GetHeader("user_id")

	ctx.Next()
	log.Printf("trace_id %s user_id %s %s begin time %d use time %d ns\n", traceId, userId, ctx.Request.RequestURI, begin.UnixNano(), time.Since(begin).Nanoseconds())
}

func main() {
	router := gin.Default()

	router.Use(traceMiddleware)
	router.GET("/greet", greet)

	router.Run("localhost:5678")
}

// go run .\micro_service\trace\self_trace\rpc_tracer\server\
