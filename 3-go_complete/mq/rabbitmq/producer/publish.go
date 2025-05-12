package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ExchangeName1 = "excg1"
	ExchangeName2 = "excg2"
	ExchangeName3 = "excg3"
)

func main2() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", user, pass))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	//创建Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()

	//声明Exchange。如果Exchange不存在会创建它；如果Exchange已存在，Server会检查声明的参数和Exchange的真实参数是否一致。
	err = ch.ExchangeDeclare(
		ExchangeName1, //Exchange Name。不要以"amq."开头，这是RabbitMQ内部使用的交换机名称。
		"fanout",      //type
		true,          //durable。Server重启后该Exchange是否还存在
		false,         //auto delete。所有bindings都退出后，是否删除该Exchange
		false,         //internal。如果设为true，则不支持publish，这种Exchange不会暴露给broker的终端用户。
		false,         //no-wait。不需要等待Server端的确认响应
		nil,           //arguments
	)
	if err != nil {
		log.Panicf("declare exchange failed: %s", err)
	}

	//发送消息。如果binding还没建好，则消息会丢失，所以发消息之前先确保subscriber已经启动好了
	produce("hello 大乔乔", ch, ExchangeName1, "")
	produce("hello world", ch, ExchangeName1, "")
	produce("hello golang", ch, ExchangeName1, "")
}

// go run .\mq\rabbitmq\producer\
