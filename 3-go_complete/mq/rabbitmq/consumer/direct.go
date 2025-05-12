package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	//guest账号只能用于连接localhost
	// user = "guest"
	// pass = "guest"
	//可以通过网页后台创建一个admin账号
	user = "dqq"
	pass = "123456"

	QueueName = "hello"
)

func createDelivery(ch *amqp.Channel, qName string) <-chan amqp.Delivery {
	deliveryCh, err := ch.Consume(
		qName, //queue
		"",    //consumer
		false, //auto-ack。autoAck其实就是noAck，只要server把消息传给consumer，本消息就会被标记为ack，而不管它有没有被consumer成功消费。
		false, //exclusive
		false, //no-local
		false, //no-wait
		nil,   //args
	)
	if err != nil {
		log.Panicf("regist consumer failed: %s", err)
	}
	return deliveryCh
}

func consumeDelivery(deliveryCh <-chan amqp.Delivery, flag int) {
	for delivery := range deliveryCh {
		log.Printf("%d receive message [%s][%s]", flag, delivery.RoutingKey, delivery.Body)
		delivery.Ack(false) //通知Server此消息已成功消费。Ack参数为true时，此channel里之前未ack的消息会一并被ack（相当于批量ack）。如果没有ack，则下一次启动时还消费到此消息（除非超时30分钟，因为delivery在30分钟后会被强制ack）,因为channel close时，它里没有ack的消息会再次被放入队列的尾部。
		// os.Exit(0)          //后面还有一些分配给该consumer的消息，会丢失
	}
}

func main() {
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

	//RoundRobin不一定是好的负载均衡方式，因为有的消息消费起来需要更多的时间。
	err = ch.Qos( //quality of service
		1,     //prefetch count。一个消费方最多能有多少条未ack的消息。如果consumer是以noAck（autoAck）启动的则server会忽略该参数。<=0时忽略该参数。该值越小，负载越均衡，但是单个消费方的吞吐也越低。
		0,     //prefetch size。Server端至少攒够这么多字节，才发给消费方。<=0时忽略该参数。
		false, //global
	)
	if err != nil {
		log.Panicf("set Qos  failed: %s", err)
	}

	//一个Connection上可以创建多个channel
	ch2, err := conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch2.Close()

	log.Printf("waiting for messages, to exit press CTRL+C")
	//一个队列可以对应多个channel（它们平分这个queue里的数据），一个channel可以有多个consumer（它们平分这个channel里的数据）。broker默认会按轮流(RoundRobin)的方式把各个消息发给所有consumer
	delivery1 := createDelivery(ch, QueueName)
	delivery2 := createDelivery(ch, QueueName)
	delivery3 := createDelivery(ch2, QueueName)
	//并行去消费消息
	go consumeDelivery(delivery1, 1)
	go consumeDelivery(delivery2, 2)
	go consumeDelivery(delivery3, 3) //注释掉试试，有消息没被打印出来，程序结束时这些消息会再次被放回队列（因为没有ack）
	select {}
}

// go run .\mq\rabbitmq\consumer\
