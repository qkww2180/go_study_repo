package redis_class

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// value是简单的string
func stringValue(ctx context.Context, client *redis.Client) {
	key := "name"
	value := "大乔乔"
	defer client.Del(ctx, key) //函数结束时删除redis上的key，不影响下次运行演示

	err := client.Set(ctx, key, value, 1*time.Second).Err() //1秒后失效。0表示永不失效
	checkError(err)

	client.Expire(ctx, key, 3*time.Second) //通过Expire设置3秒后失效。该方法对任意类型的redis value都适用
	time.Sleep(2 * time.Second)

	v2, err := client.Get(ctx, key).Result()
	checkError(err)
	fmt.Println(v2)
}

// value是List
func listValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []interface{}{1, "中", 3, 4}
	err := client.RPush(ctx, key, values...).Err() //RPush向List右侧插入，LPush向List左侧插入。如果List不存在会先创建
	checkError(err)

	v2, err := client.LRange(ctx, key, 0, -1).Result() //截取，双闭区间。LRange表示List Range，即遍历List。0表示第一个，-1表示倒数第一个。v2是个[]string，即1,3,4存到redis里实际上是string
	checkError(err)
	fmt.Println(v2)
}

// value是Set
func setValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []interface{}{1, "中", 3, 4}         //1,3,4存到redis里实际上是string
	err := client.SAdd(ctx, key, values...).Err() //SAdd向Set中添加元素
	checkError(err)

	//判断Set中是否包含指定元素
	var value any
	value = 1 //数字1会转成string再去redis里查找
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}
	value = "1"
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}
	value = 2
	if client.SIsMember(ctx, key, value).Val() {
		fmt.Printf("Set中包含%#v\n", value)
	} else {
		fmt.Printf("Set中不包含%#v\n", value)
	}

	//遍历Set
	for _, ele := range client.SMembers(ctx, key).Val() {
		fmt.Println(ele)
	}

	key2 := "ids2"
	defer client.Del(ctx, key2)
	values = []interface{}{1, "中", "大", "乔"}
	err = client.SAdd(ctx, key2, values...).Err() //SAdd向Set中添加元素
	checkError(err)

	//差集
	fmt.Println("key - key2 差集")
	for _, ele := range client.SDiff(ctx, key, key2).Val() {
		fmt.Println(ele)
	}
	fmt.Println("key2 - key 差集")
	for _, ele := range client.SDiff(ctx, key2, key).Val() {
		fmt.Println(ele)
	}

	//交集
	fmt.Println("key & key2 交集")
	for _, ele := range client.SInter(ctx, key, key2).Val() {
		fmt.Println(ele)
	}
}

// value是ZSet(有序的Set)
func zsetValue(ctx context.Context, client *redis.Client) {
	key := "ids"
	defer client.Del(ctx, key)

	values := []redis.Z{{Member: "张三", Score: 70.0}, {Member: "李四", Score: 100.0}, {Member: "王五", Score: 80.0}} //Score是用来排序的
	err := client.ZAdd(ctx, key, values...).Err()
	checkError(err)

	//遍历ZSet，按Score有序输出Member
	for _, ele := range client.ZRange(ctx, key, 0, -1).Val() {
		fmt.Println(ele)
	}
}

// value是哈希表(即map)
func hashtableValue(ctx context.Context, client *redis.Client) {
	student1 := map[string]interface{}{"Name": "张三", "Age": 18, "Height": 173.5}
	err := client.HMSet(ctx, "学生1", student1).Err() //前缀H表示HashTable。redis-server4.0之后的版本可以直接使用HSet
	checkError(err)
	student2 := map[string]interface{}{"Name": "李四", "Age": 20, "Height": 180.0}
	err = client.HMSet(ctx, "学生2", student2).Err()
	checkError(err)

	age, err := client.HGet(ctx, "学生2", "Age").Result() //指定redis的key以及map里的key
	checkError(err)
	fmt.Println(age)

	for field, value := range client.HGetAll(ctx, "学生1").Val() { //GetAll表示获取完整的map
		fmt.Println(field, value)
	}

	client.Del(ctx, "学生1")
	client.Del(ctx, "学生2")
}

func checkError(err error) {
	if err != nil {
		if err == redis.Nil { //读redis发生error，大部分情况是因为key不存在
			fmt.Println("key不存在")
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
