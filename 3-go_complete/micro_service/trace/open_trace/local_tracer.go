package main

import (
	"context"
	"dqq/micro_service/trace"
	"dqq/micro_service/util"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	trace_log "github.com/opentracing/opentracing-go/log"
)

/**
追踪调用链路上每一步的耗时，以及各个函数的调用关系
*/

func main() {
	// 链路追踪数据的记录和上报是通过tracer完成的，opentracing作为一个规范，它只提供了tracer接口，没有提供具体的tracer实现
	// opentracing.SetGlobalTracer(opentracing.NoopTracer{}) //设置一个全局的tracer。NoopTracer是一个空的tracer，仅用于mock测试或demo演示

	// jaeger实现了opentracing.tracer接口
	jaeger, closer, err := trace.NewJaegerTracer("my_service", "127.0.0.1:6831") //需要先启动jaeger
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(jaeger) //设置tracer全局 节省到处传递tracer的过程

	ctx := context.Background()
	userId := 8
	userName := "大乔乔"
	content := visitWebSite(ctx, userId, userName) //用户访问网站，函数入口
	fmt.Println(content)
}

// 用户访问网站，函数入口
func visitWebSite(ctx context.Context, userId int, userName string) string {
	// root span
	span := opentracing.GlobalTracer().StartSpan("visit_website")
	defer span.Finish() //匆必要调用Finish()，因为Finish的时候SpanConext会被写入jeager的数据库

	time.Sleep(1 * time.Millisecond)
	/**
	Tag、Log、BaggageItem都是key-value形式的数据。
	*/
	// Tag侧重于查询
	span.SetTag("访问时段", "上午")
	// Log会记录event发生的时间
	span.LogFields(
		trace_log.Int("user_visit", userId),
		trace_log.String("visit_page", "/home"),
	)
	// BaggageItem可看成是一种特殊的Log(event=baggage)，它里面的数据可以传递给子孙后代，普通Log和Tag传不给后代
	span.SetBaggageItem("trace_id", util.RandStringRunes(10))
	span.SetBaggageItem("user_id", strconv.Itoa(userId))

	recordUV(span, userId)
	// go recordUV(span, userId)
	reccommend := getReccommend(span, userId)
	return reccommend
}

// 上报用户来访，用于后续统计UV(user visit)
func recordUV(parentSpan opentracing.Span, userId int) {
	// record_uv  follows from visit_website。FollowsFrom表示parent span不以任何形式依赖child span的结果，当然child span的工作也是由parent span引起的
	span := opentracing.GlobalTracer().StartSpan("record_uv", opentracing.FollowsFrom(parentSpan.Context()))
	defer span.Finish()

	time.Sleep(3 * time.Millisecond)
}

// 调推荐的微服务，获取推荐列表
func getReccommend(parentSpan opentracing.Span, userId int) string {
	// ChildOf表示parent span某种程度上依赖child span的结果
	span := opentracing.GlobalTracer().StartSpan("reccommend", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	userRole := getUserRole(span, userId)
	list := make([]string, 0, 10)
	if "vip" != strings.ToLower(userRole) {
		list = append(list, "广告视频")
	}
	list = append(list, "gorm教程")
	list = append(list, "grpc教程")
	return strings.Join(list, "\n")
}

// 从MySQL里获取用户的角色
func getUserRole(parentSpan opentracing.Span, userId int) string {
	// ChildOf表示parent span某种程度上依赖child span的结果
	span := opentracing.GlobalTracer().StartSpan("get_user_role", opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	// visit_website --> reccommend --> get_user_role。这里的BaggageItem是从它爷爷那儿继承过来的
	fmt.Println("打印BaggageItem")
	span.Context().ForeachBaggageItem(func(k, v string) bool {
		fmt.Println(k, v)
		return true
	})
	fmt.Println()

	return "VIP"
}

// go run .\micro_service\trace\open_trace\
