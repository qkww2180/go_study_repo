package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main4() {
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
		ExchangeName3,
		"topic", //type
		true,    //durable
		false,   //auto delete
		false,   //internal
		false,   //no-wait
		nil,     //arguments
	)
	if err != nil {
		log.Panicf("declare exchange failed: %s", err)
	}

	produce("hello 大乔乔", ch, ExchangeName3, "machine1.debug")
	produce("hello world", ch, ExchangeName3, "machine2.debug")
	produce("hello golang", ch, ExchangeName3, "machine1.info")
	produce("hello dqq", ch, ExchangeName3, "machine2.info")
}

// go run .\mq\rabbitmq\producer\
