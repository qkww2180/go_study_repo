package main

import (
	"context"
	"fmt"
	"time"

	common "dqq/mq/kafka"

	"github.com/segmentio/kafka-go"
)

// 生产消息
func writeKafka(ctx context.Context) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"), //不定长参数，支持传入多个broker的ip:port
		Topic:                  common.TOPIC,                //为所有message指定统一的topic。如果这里不指定统一的Topic，则创建kafka.Message{}时需要分别指定Topic
		Balancer:               &kafka.Hash{},               //把message的key进行hash，确定partition
		WriteTimeout:           1 * time.Second,             //设定写超时
		RequiredAcks:           kafka.RequireNone,           //RequireNone不需要等待ack返回，效率最高，安全性最低；RequireOne只需要确保Leader写入成功就可以发送下一条消息；RequiredAcks需要确保Leader和所有Follower都写入成功才可以发送下一条消息。
		AllowAutoTopicCreation: true,                        //Topic不存在时自动创建。生产环境中一般设为false，由运维管理员创建Topic并配置partition数目
	}
	defer writer.Close() //记得关闭连接

	for { //无限循环，不停地往kafka里发送消息
		for i := 0; i < 3; i++ { //允许重试3次
			if err := writer.WriteMessages(ctx, //批量写入消息，原子操作，要么全写成功，要么全写失败
				kafka.Message{Key: []byte{1}, Value: []byte("大")},
				kafka.Message{Key: []byte{1}, Value: []byte("乔")},
				kafka.Message{Key: []byte{1}, Value: []byte("乔")},
				kafka.Message{Key: []byte{2}, Value: []byte("教")}, //key相同时肯定写入同一个partition
				kafka.Message{Key: []byte{2}, Value: []byte("育")},
			); err != nil {
				if err == kafka.LeaderNotAvailable { //首次写一个新的Topic时，会发生LeaderNotAvailable错误，重试一次就好了。选主操作是在第一次写数据时触发的。
					time.Sleep(500 * time.Millisecond)
					continue
				} else {
					fmt.Printf("batch write message failed: %v", err)
				}
			} else {
				break //只要成功一次就不再尝试下一次了
			}
		}
		time.Sleep(time.Second) //每发一条，稍微休息一下
	}
}

func main() {
	ctx := context.Background()
	go writeKafka(ctx)

	select {} //main永不退出
}

// go run .\mq\kafka\producer\
