package main

import (
	"dqq/concurrency/database"
	"dqq/concurrency/handler"
	"dqq/concurrency/util"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof" //在包的init()里注册了几个路由。在浏览器中访问http://127.0.0.1:8080/debug/pprof

	"github.com/gin-gonic/gin"
)

var (
	queue            = flag.String("queue", "channel", "使用哪种数据队列，channel或者kafka")
	writeOrderFinish = make(chan struct{})
)

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15。Ctrl+C对应SIGINT信号

	sig := <-c //阻塞，直到信号的到来
	util.LogRus.Infof("receive signal %s", sig.String())
	if os.Getenv("queue") == "channel" {
		handler.CloseChannel() //关闭channel
	} else {
		go handler.CloseMQ() //关闭写mq的连接
	}
	<-writeOrderFinish
	util.LogRus.Info("writeOrderFinish, exit")
	os.Exit(0) //进程退出
}

func Init() {
	flag.Parse()
	if *queue != "kafka" {
		*queue = "channel" //默认为channel
	}
	os.Setenv("queue", *queue) //写入环境变量
	util.InitLog("log")
	database.InitGiftInventory()
	if err := database.ClearOrders(); err != nil {
		panic(err)
	} else {
		util.LogRus.Info("clear table orders")
	}

	if *queue == "channel" {
		handler.InitChannel() //使用channel
	}
	go func() {
		handler.TakeOrder() //把channel里的订单信息写入Mysql
		writeOrderFinish <- struct{}{}
	}()
	go listenSignal()
	go http.ListenAndServe("127.0.0.1:8080", nil) //准备接收http://127.0.0.1:8080/debug/pprof上的访问，或者go tool pprof -http="127.0.0.1:8089"  http://127.0.0.1:8080/debug/pprof/profile

	if *queue == "kafka" {
		handler.InitMQ() //使用MQ
	}
}

func main() {
	Init()

	//GIN自带logger和recover中间件
	//[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached

	gin.SetMode(gin.ReleaseMode) //GIN线上发布模式
	// gin.DefaultWriter = io.Discard //禁止GIN的输出
	// 修改静态资源不需要重启GIN，刷新页面即可
	router := gin.Default()

	router.Static("/js", "views/js") //在url是访问目录/js相当于访问文件系统中的views/js目录
	router.Static("/img", "views/img")
	router.StaticFile("/favicon.ico", "views/img/dqq.png") //在url中访问文件/favicon.ico，相当于访问文件系统中的views/img/dqq.png文件
	router.LoadHTMLFiles("views/lottery.html")             //使用这些.html文件时就不需要加路径了

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "lottery.html", nil)
	})
	router.GET("/gifts", handler.GetAllGifts) //获取所有奖品信息
	router.GET("/lucky", handler.Lottery)     //点击抽奖按钮

	router.Run("localhost:5678")
}

// go run ./main.go -queue=channel
