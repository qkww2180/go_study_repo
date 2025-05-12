package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

// 发送消息
func produce(msg string, ch *amqp.Channel, exg, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		exg,   //exchange。""为默认的Exchange（direct类型），这种Exchange会把消息传递给routing key指定的Queue。
		key,   //routing key。Exchange为""时，routing key就是QueueName
		false, //mandatory
		false, //immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, //消息如果想持久化进磁盘，即确保RabbitMQ Server（或称broker）重启后消息不丢失，需同时满足2个条件：队列需要是durable，消息需要是Persistent。Transient显然意味着更高的吞吐。另外即使设置了Persistent，消息也不是立即会写入磁盘，中间有缓冲，如果broker突然挂掉，缓冲里的数据会丢失。
			ContentType:  "text/plain",    //MIME content type
			Body:         []byte(msg),
		},
	)
	if err != nil {
		log.Panicf("publish message failed: %s", err)
	}
	log.Printf("send [%s][%s] success", key, msg)
}

func main() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:5672/", user, pass))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	//创建Channel
	ch, err := conn.Channel() //后台会显示channels
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()

	//声明队列（队列不存在时会创建），队列用于存储消息（充当producer和consumer的中转站）。生产方和消费方都可以声明队列，如果一方声明好了，另一方就不用声明了。重复声明不会出错。如果声明的队列跟已经存在的同名队列属性不一致则会出错。
	//队列默认绑定到Name为""的Exchange，该Exchange不需要创建默认已存在，且类型为direct。
	//事实上，producer根本感知不到queue的存在，producer只是把消息发送给了Exchange，由Exchange把这条消息推送给一个或多个Queue。
	_, err = ch.QueueDeclare(
		QueueName, //队列名。为空时Server指定一个随机（且唯一）的队列名
		true,      // durable。durable=true即队列数据会持久化入磁盘，这样RabbitMQ Server重启后，队列数据也不会丢失。
		false,     // delete when unused。若设为true，当所有consumer都退出后队列会自动被删除(这意味着Server要监听Consumer进程是否存活)。
		false,     // exclusive。若设为true，队列由声明它的connection独享，其他连接上的channel不能再使用（包含声明、绑定、消费等等）同名的队列。
		false,     // no-wait。
		nil,       // arguments。
	)
	if err != nil {
		log.Panicf("declare queue failed: %s", err)
	}
	ch2, _ := conn.Channel()
	//发送消息。通过2个channel往同一个队列里面写入消息
	produce("hello 大乔乔", ch, "", QueueName)
	produce("hello world", ch, "", QueueName)
	produce("hello golang", ch2, "", QueueName)
}

// go run .\mq\rabbitmq\producer\
