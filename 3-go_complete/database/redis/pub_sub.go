package redis_class

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func publish(ctx context.Context, client *redis.Client, channel string, message any) {
	cmd := client.Publish(ctx, channel, message)
	if cmd.Err() == nil {
		n := cmd.Val() //订阅者数量
		fmt.Printf("%s向频道%s发布了消息，此时该频道有%d个订阅者\n", ctx.Value("publisher_name"), channel, n)
	} else {
		fmt.Printf("%s向频道%s发布消息失败%v\n", ctx.Value("publisher_name"), channel, cmd.Err())
	}
}

func subscribe(ctx context.Context, client *redis.Client, channels []string) {
	ps := client.Subscribe(ctx, channels...)
	defer ps.Close()

	// 无限循环，频道(channel)里一有消息，立即进行读取
	for {
		if msg, err := ps.ReceiveMessage(ctx); err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Printf("%s从频道%s里接收到消息:%s\n", ctx.Value("subscriber_name"), msg.Channel, msg.Payload)
		}
	}
}

func pubSub(ctx context.Context, client *redis.Client) {
	ctx1 := context.WithValue(ctx, "publisher_name", "publisher1")
	ctx2 := context.WithValue(ctx, "publisher_name", "publisher2")
	channel1 := "channel1"
	channel2 := "channel2"

	//启动第一批subscriber
	//subscriber要先启动好，在此之前频道(channel)里的消息它接收不到
	ctx3 := context.WithValue(ctx, "subscriber_name", "subscriber1")
	ctx4 := context.WithValue(ctx, "subscriber_name", "subscriber2")

	go subscribe(ctx3, client, []string{channel1})
	go subscribe(ctx4, client, []string{channel2})
	time.Sleep(1 * time.Second)

	go publish(ctx1, client, channel1, "白日依山尽")
	go publish(ctx2, client, channel1, "黄河入海流")
	time.Sleep(1 * time.Second)

	fmt.Println()
	//启动第二批subscriber
	ctx5 := context.WithValue(ctx, "subscriber_name", "subscriber3")
	go subscribe(ctx5, client, []string{channel1, channel2})
	time.Sleep(1 * time.Second)

	go publish(ctx1, client, channel2, "yu qiong qian li mu")
	go publish(ctx2, client, channel2, "geng shang yi ceng lou")
	time.Sleep(1 * time.Second)
}
