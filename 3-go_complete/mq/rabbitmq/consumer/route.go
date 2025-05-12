package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribeByKey(conn *amqp.Connection, flag int, exchange string, keys ...string) {
	//创建Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()
	//声明队列
	q, err := ch.QueueDeclare(
		"",    //队列名为空时Server指定一个随机（且唯一）的队列名
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("declare queue failed: %s", err)
	}

	//队列和Exchange建立绑定关系。
	//队列默认绑定到Name为""的Exchange，该Exchange不需要创建默认已存在，且类型为direct。
	for _, key := range keys {
		err = ch.QueueBind(
			q.Name, //Queue Name
			key,    //routing key。匹配上这个key的消息会被发送到这个队列
			exchange,
			false, //noWait
			nil,   //arguments
		)
		if err != nil {
			log.Panicf("bind queue failed: %s", err)
		}
	}

	//一旦开始消费，就不要再修改绑定关系了
	delivery := createDelivery(ch, q.Name)
	consumeDelivery(delivery, flag)
}

func main3() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", user, pass))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	log.Printf("waiting for messages, to exit press CTRL+C")
	go subscribeByKey(conn, 1, ExchangeName2, "debug", "info")
	go subscribeByKey(conn, 2, ExchangeName2, "error")
	go subscribeByKey(conn, 3, ExchangeName2, "debug", "info", "error")
	select {}

}

// go run .\mq\rabbitmq\consumer\
