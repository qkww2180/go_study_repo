package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Request struct {
	A int
	B int
}

type Response struct {
	Sum int
}

var (
	RequestQueue  = "sum_request"
	ResponseQueue = "sum_response"
)

var (
	requestBuffer sync.Map
	ch            *amqp.Channel
)

// 用随机字符串生成requestId
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

func receiveResponse() {
	deliveryCh := createDelivery(ch, ResponseQueue)
	for delivery := range deliveryCh { //从MQ里取出response
		delivery.Ack(false)
		requestId := delivery.CorrelationId
		if v, exists := requestBuffer.Load(requestId); exists { //从callBuffer里取出requestId对应的channel
			requestBuffer.Delete(requestId) //从buffer里删除requestId
			resp := string(delivery.Body)
			ch := v.(chan string)
			ch <- resp //把响应放到channel里
		} else {
			log.Printf("discard response %s", requestId) //如果server端重复消费request，client端就会重复收到response，第二次收到response时相应的requestId已经不存在了，即忽略该response
		}
	}
}

// 远程过程调用
func call(a, b int) int {
	request := Request{A: a, B: b}
	bs, err := sonic.Marshal(request)
	if err != nil {
		log.Printf("marshal request failed: %s", err)
		return 0
	}
	requestId := randomString(10) //requestId是一个长度为10的随机字符串
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",           //exchange
		RequestQueue, // routing key，即Queue Name。请求发到这个队列中
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:   "application/json", //MIME
			CorrelationId: requestId,
			ReplyTo:       ResponseQueue, //希望响应发到这个队列中
			Body:          bs,
		},
	)
	if err != nil {
		log.Printf("send request failed: %s", err)
		return 0
	}
	respCh := make(chan string, 1)
	requestBuffer.Store(requestId, respCh)
	s := <-respCh //阻塞，直到能从channel里取出响应
	var resp Response
	err = sonic.Unmarshal([]byte(s), &resp)
	if err != nil {
		log.Printf("unmarshal response failed: %s", err)
		return 0
	} else {
		return resp.Sum
	}
}

func main5() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", user, pass))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	//创建Channel
	ch, err = conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()

	//声明队列
	_, err = ch.QueueDeclare(
		ResponseQueue, //queue name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Panicf("declare queue failed: %s", err)
	}

	//从消息队列里接收response，放入callBuffer
	go receiveResponse()

	requestCnt := 10 //并发发起多次请求
	wg := sync.WaitGroup{}
	wg.Add(requestCnt)
	for i := 0; i < requestCnt; i++ {
		go func() {
			defer wg.Done()
			a := rand.Intn(10)
			b := rand.Intn(10)
			c := call(a, b)
			log.Printf("%d + %d = %d", a, b, c) //打印结果，看看并发RPC能否正常工作
		}()
	}
	wg.Wait()
}

// 可以启多个Client进程
// go run .\mq\rabbitmq\consumer\
