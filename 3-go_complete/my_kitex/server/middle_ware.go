package main

import (
	"context"
	"errors"
	"fmt"
	math_service "my_kitex/kitex_gen/math_service"
	"time"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/pkg/utils/kitexutil"
	"google.golang.org/grpc/metadata"
)

/*
server端支持设置多种中间件，它们的调用顺序为：
1、server.WithMiddleware 设置的中间件，按其在 Option 中的设置顺序执行
2、ACLMiddleware，参见自定义访问控制
3、server.WithErrorHandler 设置的中间件
*/

// 计时中间件
func TimerMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		begin := time.Now()
		err := next(ctx, request, response)

		//获取框架RPC信息
		caller, _ := kitexutil.GetCaller(ctx)   //获取调用方的ServiceName
		addr, _ := kitexutil.GetCallerAddr(ctx) //获取调用方地址
		method, _ := kitexutil.GetMethod(ctx)

		klog.Infof("%s %s %s use time %d ms", caller, addr, method, time.Since(begin).Milliseconds())
		return err
	}
}

// 打印参数中间件
func RecordArgMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		//打印请求参数
		if arg, ok := request.(utils.KitexArgs); ok {
			switch req := arg.GetFirstArgument().(type) {
			case *math_service.AddRequest:
				klog.Debugf("AddRequest Left %d Right %d", req.Left, req.Right)
			case *math_service.SubRequest:
				klog.Debugf("SubRequest Left %d Right %d", req.Left, req.Right)
			default:
				klog.Debug("Request %v", req)
			}
		}
		//执行下一个中间件
		err := next(ctx, request, response)
		ri := rpcinfo.GetRPCInfo(ctx) //获取框架RPC信息
		if stats := ri.Stats(); stats != nil {
			//打印panic信息
			if panicHappened, panicInfo := stats.Panicked(); panicHappened { //panicInfo就是recover()的返回值
				klog.Errorf("panic info %s", panicInfo)
				return fmt.Errorf("服务内部发生panic") //原始panic信息对client隐藏
			} else {
				if stats.Error() == nil {
					//打印响应参数
					if result, ok := response.(utils.KitexResult); ok {
						switch resp := result.GetResult().(type) {
						case *math_service.AddResponse:
							if resp == nil {
								klog.Error("kitex result is nil")
							} else {
								klog.Debugf("AddResponse Sum %d", resp.Sum)
							}
						case *math_service.SubResponse:
							if resp == nil {
								klog.Error("kitex result is nil")
							} else {
								klog.Debugf("SubResponse Diff %d", resp.Diff)
							}
						default:
							klog.Debug("Response %v", resp)
						}

					}
				} else {
					klog.Errorf("biz error: %s", stats.Error().Error())
				}
			}
		}
		return err
	}
}

func AuthMW(ctx context.Context, request any) (reason error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok { //BUG: 在kitex中metadata数据被丢弃了，取不到
		if value, exists := md["token"]; exists {
			klog.Debugf("token is %s", value[0])
			if value[0] == "123456" {
				return nil //返回nil表示认证通过
			}
		} else {
			klog.Error("token not found in metadata")
		}
	} else {
		klog.Error("not found incoming context")
	}
	return errors.New("invalid token") //返回error表示认证不通过
}
