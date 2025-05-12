package common

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 向Jaeger上报数据。删掉这个中间件不影响业务的正常运行
func ServerTraceMiddleware(ctx *gin.Context) {
	//拿着Uber-Trace-Id去查询jaeger数据库，查询结果反序列化为一个SpanContext
	clientSpanCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(ctx.Request.Header))
	if err != nil {
		log.Printf("反序列化request header失败: %s", err)
	}

	// 创建server端span
	operationName := ctx.Request.RequestURI
	serverSpan := opentracing.StartSpan(
		operationName,
		ext.RPCServerOption(clientSpanCtx)) //把client的metadata带进来。一方面指明继承关系，另一方面指明这是一个Server端的span(Tag里有一条span.kind=server)
	defer serverSpan.Finish() //将SpanContext写入jaeger数据库
	for k, v := range ctx.Request.Header {
		if k == "Uber-Trace-Id" {
			continue
		}
		serverSpan.SetTag(k, v[0])
	}

	ctx.Next()
}
