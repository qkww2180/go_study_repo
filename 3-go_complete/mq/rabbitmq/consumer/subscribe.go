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

func subscribe(ch *amqp.Channel, flag int, exchange string) {
	//声明队列
	q, err := ch.QueueDeclare(
		"",    //队列名为空时Server指定一个随机（且唯一）的队列名
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("declare queue failed: %s", err)
	}

	//队列和Exchange建立绑定关系
	err = ch.QueueBind(
		q.Name, //Queue Name
		"",     //routing key。fout模式下会忽略routing key
		exchange,
		false, //noWait
		nil,   //arguments
	)
	if err != nil {
		log.Panicf("bind queue failed: %s", err)
	}

	delivery := createDelivery(ch, q.Name)
	consumeDelivery(delivery, flag)
}

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

	log.Printf("waiting for messages, to exit press CTRL+C")
	go subscribe(ch, 1, ExchangeName1)
	go subscribe(ch, 2, ExchangeName1)
	go subscribe(ch, 3, ExchangeName1)
	select {}
}

// go run .\mq\rabbitmq\consumer\
