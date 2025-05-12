package main

import (
	"context"
	"dqq/micro_service/util"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

/**
追踪调用链路上每一步的耗时
*/

func main() {
	ctx := context.Background()
	userId := 8
	userName := "大乔乔"
	ctx = context.WithValue(ctx, "trace_id", util.RandStringRunes(20))
	content := visitWebSite(ctx, userId, userName) //用户访问网站，函数入口
	fmt.Println(content)
}

// 用户访问网站，函数入口
func visitWebSite(ctx context.Context, userId int, userName string) string {
	begin := time.Now()
	defer func() {
		funcName, _, _, _ := runtime.Caller(0)
		fn := runtime.FuncForPC(funcName).Name()
		log.Printf("trace_id %s %s use time %d ns", ctx.Value("trace_id").(string), fn, time.Since(begin).Nanoseconds())
	}()
	go recordUV(ctx, userId)
	reccommend := getReccommend(ctx, userId)
	return reccommend
}

// 上报用户来访，用于后续统计UV(user visit)
func recordUV(ctx context.Context, userId int) {
	begin := time.Now()
	defer func() {
		funcName, _, _, _ := runtime.Caller(0)
		fn := runtime.FuncForPC(funcName).Name()
		log.Printf("trace_id %s %s use time %d ns", ctx.Value("trace_id").(string), fn, time.Since(begin).Nanoseconds())
	}()
}

// 从MySQL里获取用户的角色
func getUserRole(ctx context.Context, userId int) string {
	begin := time.Now()
	defer func() {
		funcName, _, _, _ := runtime.Caller(0)
		fn := runtime.FuncForPC(funcName).Name()
		log.Printf("trace_id %s %s use time %d ns", ctx.Value("trace_id").(string), fn, time.Since(begin).Nanoseconds())
	}()
	return "VIP"
}

// 调推荐的微服务，获取推荐列表
func getReccommend(ctx context.Context, userId int) string {
	begin := time.Now()
	defer func() {
		funcName, _, _, _ := runtime.Caller(0)
		fn := runtime.FuncForPC(funcName).Name()
		log.Printf("trace_id %s %s use time %d ns", ctx.Value("trace_id").(string), fn, time.Since(begin).Nanoseconds())
	}()
	userRole := getUserRole(ctx, userId)
	list := make([]string, 0, 10)
	if "vip" != strings.ToLower(userRole) {
		list = append(list, "广告视频")
	}
	list = append(list, "gorm教程")
	list = append(list, "grpc教程")
	return strings.Join(list, "\n")
}

// go run .\micro_service\trace\self_trace\
