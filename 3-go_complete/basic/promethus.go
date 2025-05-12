package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	temperature = promauto.NewGauge(prometheus.GaugeOpts{Name: "temperature"})
)

func main19() {
	//启动一个http server，在/metrics路径上处理promethus的抓取请求
	go func() {
		http.Handle("/metrics", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				promhttp.Handler().ServeHTTP(w, r)
			}))

		if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
			fmt.Println(err)
		}
	}()

	//业务主流程
	for {
		temperature.Set(rand.Float64() * 100) // 每隔1秒钟上报一次CPU的温度
		time.Sleep(1 * time.Second)
	}
}
