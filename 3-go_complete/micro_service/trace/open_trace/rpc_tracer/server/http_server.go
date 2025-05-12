package main

import (
	"dqq/micro_service/trace"
	"dqq/micro_service/trace/open_trace/rpc_tracer/common"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func greet(ctx *gin.Context) {
	time.Sleep(100 * time.Millisecond)
	ctx.String(http.StatusOK, "hello")
}

var (
	closer io.Closer
)

func InitHttpJaeger() {
	var jaeger opentracing.Tracer
	var err error
	// 设置全局tracer
	jaeger, closer, err = trace.NewJaegerTracer("my_http_server", "127.0.0.1:6831") //需要先启动jaeger
	if err != nil {
		panic(err)
	}
	// defer closer.Close()  需要在接收到kill信号时调用Close()
	opentracing.SetGlobalTracer(jaeger)
}

func main() {
	InitHttpJaeger()
	router := gin.Default()

	router.Use(common.ServerTraceMiddleware)
	router.GET("/greet", greet)

	router.Run("localhost:5678")
}

// go run .\micro_service\trace\open_trace\rpc_tracer\server\
