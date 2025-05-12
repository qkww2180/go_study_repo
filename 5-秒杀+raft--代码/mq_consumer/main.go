package main

import (
	"context"
	"dqq/concurrency/database"
	"dqq/concurrency/util"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/segmentio/kafka-go"
)

var reader *kafka.Reader

func Init() {
	util.InitLog("mq_consumer_log")
	viper := util.CreateConfig("kafka")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(viper.GetString("brokers"), ","),
		Topic:          viper.GetString("topic"),
		StartOffset:    kafka.LastOffset,  //之前MQ里的老数据不再接收了
		GroupID:        "serialize_order", //注意：如果不指定GroupID，则只能消费到1个partition里的数据。我的kafka配的是生成2个partition
		CommitInterval: 1 * time.Second,   //每隔多长时间自动commit一次offset
	})
	util.LogRus.Info("create reader to mq")
}

// 从mq里取出订单，把订单写入Mysql
func ConsumeOrder() {
	for {
		if message, err := reader.ReadMessage(context.Background()); err != nil {
			fmt.Printf("read message from mq failed: %v", err)
			break
		} else {
			var order database.Order
			if err := sonic.Unmarshal(message.Value, &order); err == nil {
				util.LogRus.Debugf("message partition %d", message.Partition)
				database.CreateOrder(order.UserId, order.GiftId) //写入mysql
			} else {
				util.LogRus.Errorf("order info is invalid json format: %s", string(message.Value))
			}
		}
	}
}

func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号。Ctrl+C对应SIGINT信号
	sig := <-c                                        //阻塞，直到信号的到来
	reader.Close()
	util.LogRus.Infof("receive signal %s, exit", sig.String())
	os.Exit(0) //进程退出

}

func main() {
	Init()
	go listenSignal()
	ConsumeOrder()
}

// go run ./mq_consumer/

// 启动kafka
// 开一个终端，启动zookeeper：
// D:\software\kafka_2.12-3.5.0> .\bin\windows\zookeeper-server-start.bat .\config\zookeeper.properties
// 另开一个终端，启动kafka：
// PS D:\software\kafka_2.12-3.5.0> .\bin\windows\kafka-server-start.bat .\config\server.properties
