package main

import (
	"context"
	"dqq/micro_service/trace"
	"dqq/micro_service/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
)

var (
	closer io.Closer
)

func InitHttpJaeger() {
	var jaeger opentracing.Tracer
	var err error
	// 设置全局tracer
	jaeger, closer, err = trace.NewJaegerTracer("my_http_client", "127.0.0.1:6831") //需要先启动jaeger
	if err != nil {
		panic(err)
	}
	// defer closer.Close()  需要在接收到kill信号时调用Close()
	opentracing.SetGlobalTracer(jaeger)
}

func greet(ctx context.Context) string {
	url := "http://127.0.0.1:5678/greet"
	client := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("new h_http request failed: %s", err)
		return ""
	}
	request.Header.Set("trace_id", ctx.Value("trace_id").(string))
	request.Header.Set("user_id", ctx.Value("user_id").(string))

	//创建span仅仅是为发向Jaeger上报数据，删掉这段代码不影响业务的正常运行
	span := opentracing.StartSpan("greet")
	defer span.Finish() //SpanContext写入jaeger数据库
	if span != nil {
		for k, v := range request.Header {
			span.SetTag(k, v[0])
		}
		if err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(
				request.Header,
			),
		); err != nil {
			log.Printf("序列化context失败")
		} else {
			for k, v := range request.Header { //header里会多出一条Uber-Trace-Id
				fmt.Println(k, v[0])
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		log.Printf("call h_http server failed: %s", err)
		return ""
	}
	defer response.Body.Close()
	bs, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("read h_http response failed: %s", err)
		return ""
	}
	return string(bs)
}

func main() {
	InitHttpJaeger()
	userId := 8
	ctx := context.WithValue(context.Background(), "user_id", strconv.Itoa(userId))
	ctx = context.WithValue(ctx, "trace_id", util.RandStringRunes(10))
	fmt.Println(greet(ctx))

	closer.Close()
}

// go run .\micro_service\trace\open_trace\rpc_tracer\client\
