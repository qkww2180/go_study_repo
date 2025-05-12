package main

import (
	"context"
	"time"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils/kitexutil"
	"google.golang.org/grpc/metadata"
)

/*
client端支持设置多种中间件，它们的调用顺序为：
1、xDS 路由、服务级别熔断、超时
2、client.WithContextMiddlewares 设置的中间件
3、client.WithMiddleware 设置的中间件，按其在 Option 中的设置顺序执行
4、ACLMiddleware，参见自定义访问控制
5、服务发现、实例熔断、实例级 Middleware / 服务发现、代理 Middleware
6、client.WithErrorHandler 设置的中间件
*/

// 计时中间件
func TimerMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		begin := time.Now()
		err := next(ctx, request, response)
		var service, method string
		// ri := rpcinfo.GetRPCInfo(ctx) //获取框架RPC信息
		// if ivk := ri.Invocation(); ivk != nil {
		// 	service = ivk.ServiceName()
		// 	method = ivk.MethodName()
		// }

		//获取service、method的第二种方式
		service, _ = kitexutil.GetIDLServiceName(ctx)
		method, _ = kitexutil.GetMethod(ctx)

		klog.Infof("%s %s use time %d ms", service, method, time.Since(begin).Milliseconds())
		return err
	}
}

// 基于简单口令的认证中间件
func AuthMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "token", "123456") //BUG: 在kitex中metadata数据被丢弃了，取不到
		err := next(ctx, request, response)
		return err
	}
}
