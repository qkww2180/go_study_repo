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

	log.Printf("waiting for messages, to exit press CTRL+C")
	go subscribeByKey(conn, 1, ExchangeName3, "*.info")
	go subscribeByKey(conn, 2, ExchangeName3, "machine1.*")
	go subscribeByKey(conn, 3, ExchangeName3, "*.*")
	select {}
}

// go run .\mq\rabbitmq\consumer\
