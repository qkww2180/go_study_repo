package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RequestQueue  = "sum_request"
	ResponseQueue = "sum_response"
)

type Request struct {
	A int
	B int
}

type Response struct {
	Sum int
}

// 核心接口实现
func Sum(request *Request) *Response {
	time.Sleep(time.Duration(rand.IntN(500)) * time.Millisecond) //随机休息一段时间
	return &Response{Sum: request.A + request.B}
}

type Consumer = <-chan amqp.Delivery

// 一个amqp.Channel上可以创建多个consumer
func createConsumer(ch *amqp.Channel, qName string) Consumer {
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

func main5() {
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

	//声明队列
	q, err := ch.QueueDeclare(
		RequestQueue, //queue name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Panicf("declare queue failed: %s", err)
	}

	//多台server做负载均衡
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Panicf("set Qos  failed: %s", err)
	}

	deliveryCh := createConsumer(ch, q.Name)
	const P = 10
	for i := 0; i < P; i++ { //Server端开多个协程处理请求
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			for delivery := range deliveryCh { //此for循环永不能出
				delivery.Ack(false)
				var request Request
				var response *Response
				if err := sonic.Unmarshal(delivery.Body, &request); err != nil {
					log.Printf("invalid request %s", string(delivery.Body))
					continue
				} else {
					response = Sum(&request)
				}
				if resp, err := sonic.Marshal(response); err != nil {
					log.Printf("marchal response failed %s", err)
					continue
				} else {
					err = ch.PublishWithContext(
						ctx,
						"",               //Exchange
						delivery.ReplyTo, //routing key，即Queue Name
						false,            // mandatory
						false,            // immediate
						amqp.Publishing{
							ContentType:   "application/json",     //json格式的数据
							CorrelationId: delivery.CorrelationId, //通过CorrelationId将request和response进行对应
							Body:          resp,
						},
					)
					if err != nil {
						log.Printf("send response failed: %s", err)
						continue
					} else {
						log.Printf("%d + %d = %d", request.A, request.B, response.Sum)
					}
				}
			}
		}()
	}
	log.Printf("waiting for RPC request, to exit press CTRL+C")
	select {}
}

// 可以启多个Server进程
// go run .\mq\rabbitmq\producer\
