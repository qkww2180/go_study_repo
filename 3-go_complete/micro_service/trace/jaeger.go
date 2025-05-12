package trace

import (
	"io"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func NewJaegerTracer(serviceName string, jaegerHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1, //全部采样
		},
		Reporter: &jaegercfg.ReporterConfig{
			//当span发送到服务器时要不要打日志
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jaegerHost,
		},
		ServiceName: serviceName,
	}
	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jaeger.NullLogger), //jaeger系统本身的日志输出到哪儿
		jaegercfg.Metrics(metrics.NullFactory),
	)

	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
	}
	return
}
