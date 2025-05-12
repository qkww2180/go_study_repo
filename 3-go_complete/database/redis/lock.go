package redis_class

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// TryLock 尝试获得分布式锁，成功返回true，失败返回false
func TryLock(rc *redis.Client, key string, expire time.Duration) bool {
	cmd := rc.SetNX(context.Background(), key, "value随意", expire) //SetNX如果key不存在则返回true，写入key，并设置过期时间
	if cmd.Err() != nil {
		return false
	} else {
		return cmd.Val()
	}
}

// ReleaseLock 释放分布式锁
func ReleaseLock(rc *redis.Client, key string) {
	rc.Del(context.Background(), key)
}

func LockRace(client *redis.Client) {
	key := "iPhone秒杀-共3部-第2部"      //模拟：谁抢到锁，谁就能抢到这部iPhone
	defer ReleaseLock(client, key) //函数结束时删除redis上的key，不影响下次运行演示
	const P = 100                  //100个人来抢
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func(i int) {
			defer wg.Done()
			if TryLock(client, key, time.Hour) { //秒杀活动只持续10分钟，1小时后自动把key从redis上删掉（节约redis内存），不用再调用ReleaseLock()函数
				fmt.Printf("协程%d抢到锁\n", i)
			}
		}(i)
	}
	wg.Wait()
}
