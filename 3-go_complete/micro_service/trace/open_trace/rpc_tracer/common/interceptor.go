package common

import (
	"context"
	"dqq/micro_service/util"
	"fmt"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MDCarrier完成是仿照opentracing.TextMapCarrier写的，不同之处在于TextMapCarrier是map[string]string。
//
//	type TextMapCarrier map[string]string
//
// MDCarrier同时实现了opentracing.TextMapReader和opentracing.TextMapWriter接口
type MDCarrier map[string][]string

func (m MDCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, strs := range m {
		for _, v := range strs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m MDCarrier) Set(key, val string) {
	m[key] = append(m[key], val)
}

func ServerTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//RPC对端传过来一些metadata，其中包含了uber-trace-id
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	for k, v := range md {
		fmt.Println(k, v[0]) //能看到有一条uber-trace-id
	}
	//拿着uber-trace-id去查询jaeger数据库，查询结果反序列化为一个SpanContext
	clientSpanCtx, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		MDCarrier(md),
	)
	if err != nil {
		log.Printf("Extract SpanContext failed: %s", err)
		return handler(ctx, req)
	}

	//创建server端span
	serverSpan := opentracing.GlobalTracer().StartSpan(
		info.FullMethod,                    //用完整的方法名为span命名
		ext.RPCServerOption(clientSpanCtx), //把client的metadata带进来。一方面指明继承关系，另一方面指明这是一个Server端的span(Tag里有一条span.kind=server)
		opentracing.Tag{Key: string(ext.Component), Value: "grpc服务端"}, //为span添加tag
	)
	for k, v := range md {
		if k == "uber-trace-id" {
			continue
		}
		serverSpan.SetTag(k, v[0]) //为span添加tag
	}
	defer serverSpan.Finish()

	return handler(ctx, req)
}

func ClientTraceInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	traceId := util.RandStringRunes(10)         //生成随机字符串
	md, ok := metadata.FromOutgoingContext(ctx) //OutgoingContext里可能本来就有一些数据
	if !ok {                                    //也可能没有
		md = metadata.New(map[string]string{})
	}
	md.Append("trace_id", traceId)
	md.Append("user_id", ctx.Value("user_id").(string))

	//创建span仅仅是为发向Jaeger上报数据，删掉这段代码不影响业务的正常运行
	span := opentracing.GlobalTracer().StartSpan(
		method, //以GRPC方法名为Span命名
		//为span添加tag
		opentracing.Tag{Key: string(ext.Component), Value: "grpc客户端"}, //ext.Component是一个全局变量，就等于"component"
		ext.SpanKindRPCClient, //标记为client端span。等价于加了个Tag span.kind=client
	)
	defer span.Finish() //匆必要调用Finish()，因为Finish的时候SpanConext会被写入jeager的数据库
	for k, v := range md {
		span.SetTag(k, v[0]) //为span添加tag
		fmt.Println(k, v[0]) //调Inject之前看看md里有什么
	}
	if err := opentracing.GlobalTracer().Inject( //执行Inject时，jaeger会把span里的信息注入到md里去，md里会多一条数据uber-trace-id
		span.Context(),
		opentracing.TextMap,
		MDCarrier(md),
	); err != nil {
		log.Printf("Inject SpanContext failed: %s", err)
	}
	fmt.Println("--------------------------")
	for k, v := range md {
		fmt.Println(k, v[0]) //md里会多一条数据uber-trace-id
	}
	fmt.Println("--------------------------")

	ctx = metadata.NewOutgoingContext(ctx, md) //执行RPC时需要把uber-trace-id传给对方，否则对方在执行Extract时会报错：SpanContext not found in Extract carrier
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
