package main

import (
	"context"
	"log"
	"testing"
	"time"

	redis "github.com/redis/go-redis/v9"
	etcdv3 "go.etcd.j_io/etcd/client/v3"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", //没有密码
		DB:       0,  //redis默认会创建0-15号DB，这里使用默认的DB
	})
}

func NewEtcdClient() (*etcdv3.Client, error) {
	return etcdv3.New(
		etcdv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 3 * time.Second,
		},
	)
}

const (
	KEY   = "up_name"
	VALUE = "高性能golang"
)

func BenchmarkRedisSet(b *testing.B) {
	ctx := context.Background()
	cl := NewRedisClient()
	var err error
	for n := 0; n < b.N; n++ {
		err = cl.Set(ctx, KEY, VALUE, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkRedisGet(b *testing.B) {
	ctx := context.Background()
	cl := NewRedisClient()
	var err error
	for n := 0; n < b.N; n++ {
		err = cl.Get(ctx, KEY).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkEtcdSet(b *testing.B) {
	ctx := context.Background()

	cl, err := NewEtcdClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if er := cl.Close(); er != nil {
			log.Fatal(er)
		}
	}()

	for n := 0; n < b.N; n++ {
		_, err = cl.Put(ctx, KEY, VALUE)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkEtcdGet(b *testing.B) {
	ctx := context.Background()

	cl, err := NewEtcdClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if er := cl.Close(); er != nil {
			log.Fatal(er)
		}
	}()

	for n := 0; n < b.N; n++ {
		_, err = cl.Get(ctx, KEY)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// go test -bench=Redis .\micro_service\etcd\ -benchmem
// go test -bench=Etcd .\micro_service\etcd\ -benchmem

/**
windows11上单机版最简安装Redis和Etcd。
etcd-v3.5.9-windows-amd64
Redis version 3.0.504 64 bit
启动redis命令  C:\Program Files\Redis> .\redis-server.exe .\redis.windows-service.conf
启动etcd命令   etcd

试验结果：
goos: windows
goarch: amd64
pkg: upgrading/micro_service/etcd
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkRedisSet-8        34846             30115 ns/op             200 B/op          5 allocs/op
BenchmarkRedisGet-8        35486             29403 ns/op             172 B/op          5 allocs/op
BenchmarkEtcdSet-8           313           3397384 ns/op            8687 B/op        116 allocs/op
BenchmarkEtcdGet-8          7298            185779 ns/op            7468 B/op        123 allocs/op

结论：
无论从内存开销还是速度上，Redis都完胜。Redis的读写性能相当，Etcd的写远慢于读。
*/
