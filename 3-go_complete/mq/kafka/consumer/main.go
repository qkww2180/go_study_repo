package main

import (
	"context"
	common "dqq/mq/kafka"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	reader *kafka.Reader
)

// 列出所有的partition以及它对应的topic
func listTopicAndPartition() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	for _, p := range partitions {
		if strings.HasPrefix(p.Topic, "__") {
			continue
		}
		fmt.Printf("partition %d topic %s\n", p.ID, p.Topic)
	}
}

// 消费消息
func readKafka(ctx context.Context) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"}, //支持传入多个broker的ip:port
		Topic:          common.TOPIC,
		CommitInterval: 1 * time.Second, //每隔多长时间自动commit一次offset。即一边读一边向kafka上报读到了哪个位置。
		GroupID:        "recommend_biz", //一个Group内消费到的消息不会重复。注意：如果不指定GroupID，则只能消费到1个partition里的数据，所以consumer的个数需要多于partition数据才能把数据消费全
		// Partition:      0,                 //Partition和GroupID不能同时指定
		StartOffset: kafka.FirstOffset, //当一个特定的partition没有commited offset时(比如第一次读一个partition，之前没有commit过)，通过StartOffset指定从第一个还是最后一个位置开始消费。StartOffset的取值要么是FirstOffset要么是LastOffset，LastOffset表示Consumer启动之前生成的老数据不管了。仅当指定了GroupID时，StartOffset才生效。
	})
	// defer reader.Close() //由于下面是死循环，正常情况下readKafka()函数永远不会结束，defer不会执行。所以需要监听信息2和15，当收到信号时关闭reader。需要把reader设为全局变量

	for { //消息队列里随时可能有新消息进来，所以这里是死循环，类似于读Channel
		if message, err := reader.ReadMessage(ctx); err != nil {
			fmt.Printf("read message from kafka failed: %v", err)
			break
		} else {
			offset := message.Offset
			fmt.Printf("topic=%s, partition=%d, offset=%d, key=%s, message content=%s\n", message.Topic, message.Partition, offset, string(message.Key), string(message.Value))
		}
	}
}

// 需要监听信号2和15，当收到信号时关闭reader
func listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) //注册信号2和15
	sig := <-c                                        //阻塞，直到信号的到来
	fmt.Printf("receive signal %s\n", sig.String())
	if reader != nil {
		reader.Close() //把没提交的提交掉
	}
	os.Exit(0) //进程退出
}

func main() {
	ctx := context.Background()
	listTopicAndPartition()
	go listenSignal()
	readKafka(ctx)
}

// 开3个终端，启动3个consumer试试。由于我们在kafka的config\server.properties里配的partition数目是2，所以3个consumer中有一个获取不到消息
// go run .\mq\kafka\consumer\
