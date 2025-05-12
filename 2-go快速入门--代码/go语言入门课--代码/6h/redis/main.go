package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func stringValue(ctx context.Context, client *redis.Client) {
	key := "name"
	value := "大乔乔"
	err := client.Set(ctx, key, value, 1*time.Second).Err() //1秒后失效。0表示永不失效
	checkError(err)

	client.Expire(ctx, key, 3*time.Second) //通过Expire设置3秒后失效
	time.Sleep(2 * time.Second)

	v2, err := client.Get(ctx, key).Result()
	checkError(err)
	fmt.Println(v2)

	client.Del(ctx, key)
}

func listValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	values := []interface{}{1, "中", 3, 4}
	err := client.RPush(ctx, key, values...).Err() //向List右侧插入。如果List不存在会先创建
	checkError(err)

	v2, err := client.LRange(ctx, key, 0, -1).Result() //截取。双闭区间
	checkError(err)
	fmt.Println(v2)

	client.Del(ctx, key)
}

func hashtableValue(ctx context.Context, client *redis.Client) {
	err := client.HMSet(ctx, "学生1", map[string]interface{}{"Name": "张三", "Age": 18, "Height": 173.5}).Err()
	checkError(err)
	err = client.HMSet(ctx, "学生2", map[string]interface{}{"Name": "李四", "Age": 20, "Height": 180.0}).Err()
	checkError(err)

	age, err := client.HGet(ctx, "学生2", "Age").Result()
	checkError(err)
	fmt.Println(age)

	for field, value := range client.HGetAll(ctx, "学生1").Val() {
		fmt.Println(field, value)
	}

	// client.Del(ctx, "学生1")
	// client.Del(ctx, "学生2")
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", //没有密码
		DB:       0,  //redis默认会创建0-15号DB，这里使用默认的DB
	})
	ctx := context.TODO()
	stringValue(ctx, client)
	listValue(ctx, client)
	hashtableValue(ctx, client)
}

func checkError(err error) {
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key不存在")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// 先确保启动了Redis服务  sudo service redis-server start
// go run ./redis
