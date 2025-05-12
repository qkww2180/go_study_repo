package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	/**
	Counter的计数是整数，Gauge的计数是float64。Counter可加可减，是一个历史累积量，Gauge是某一次的具体值(各次上报是独立的)
	Vec可以附带多个标签，即WithLabelValues()可接收不定长参数。不带vec则没有WithLabelValues()方法可调用
	*/
	ageGuage       = promauto.NewGauge(prometheus.GaugeOpts{Name: "user_age"})
	requestTimer   = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "request_time"}, []string{"interface"})
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{Name: "request_counter"}, []string{"interface"})
)

// 计时中间件。通用的监控上报，在中间件里做。如果是GRPC服务则在拦截器里做监控上报
func timerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()
		ifc := ctx.Request.RequestURI             //根据请求路径标识接口
		requestCounter.WithLabelValues(ifc).Inc() //接口请求次数加1。WithLabelValues()参数与NewGaugeVec()时的Label一一对应
		ctx.Next()
		requestTimer.WithLabelValues(ifc).Set(float64(time.Since(begin).Milliseconds())) //上报接口的耗时
	}
}

type User struct {
	Name string
	Age  int
}

func getUserAge(ctx *gin.Context) {
	user := User{
		Age: rand.Intn(100),
	}
	ageGuage.Set(float64(user.Age))
	ctx.JSON(http.StatusOK, user)
}

func getUserName(ctx *gin.Context) {
	user := User{
		Name: "高性能golang",
	}
	ctx.JSON(http.StatusOK, user)
}

func main() {
	engine := gin.Default()
	engine.Use(timerMiddleware())
	engine.GET("/name", getUserName)
	engine.GET("/age", getUserAge)

	//如果是GRPC服务，需要单独开一个协程启一个http server，供prometheus来取数据
	engine.GET("/metrics", func(ctx *gin.Context) { //prometheus会来这个接口上拉取监控数据
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	engine.Run("127.0.0.1:5678")
}

// go run .\micro_service\monitor

// 打开浏览器，访问
// http://127.0.0.1:5678/name
// http://127.0.0.1:5678/age
// http://127.0.0.1:5678/metrics
