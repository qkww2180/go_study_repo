package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	etcdv3 "go.etcd.j_io/etcd/client/v3"
)

const ConfigPrefix = "dqq_"

type GlobalConfig struct {
	Thresh int
	Name   string
}

var Config GlobalConfig //搞成全局变量，或者搞成单例模式

// 从Etcd上读取初始的配置参数
func InitGlobalConfig(ctx context.Context, client *etcdv3.Client) {
	if response, err := client.Get(context.Background(), ConfigPrefix+"thresh"); err == nil {
		if len(response.Kvs) > 0 {
			thresh, err := strconv.Atoi(string(response.Kvs[0].Value))
			if err == nil {
				Config.Thresh = thresh
				fmt.Printf("从etcd上获得thresh为%d\n", Config.Thresh)
			}
		}
	}
	if response, err := client.Get(context.Background(), ConfigPrefix+"name"); err == nil {
		if len(response.Kvs) > 0 {
			Config.Name = string(response.Kvs[0].Value)
			fmt.Printf("从etcd上获得name为%s\n", Config.Name)
		}
	}
}

// 时刻监听dqq_name和dqq_thresh这两个key对应的value是否发生变化
func watch(ctx context.Context, client *etcdv3.Client) {
	ch := client.Watch(ctx, ConfigPrefix, etcdv3.WithPrefix()) //监听全局配置的变化，每一个修改都会放入管道ch。WithPrefix()指明了这里的key实际上只是前缀，而非完整的key
	for response := range ch {                                 //遍历管道。这是个死循环，除非关闭管道
		for _, event := range response.Events {
			// fmt.Printf("事件类型:%s, key:%s, value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			if "PUT" == event.Type.String() { //只关心PUT事件，即数据更新
				switch string(event.Kv.Key) {
				case ConfigPrefix + "name":
					Config.Name = string(event.Kv.Value)
					fmt.Printf("name更新为%s\n", Config.Name)
				case ConfigPrefix + "thresh":
					thresh, err := strconv.Atoi(string(event.Kv.Value))
					if err == nil {
						Config.Thresh = thresh
						fmt.Printf("thresh更新为%d\n", Config.Thresh)
					}
				}
			}
		}
	}
}

func UseConfig() {
	fmt.Printf("name=%s, thresh=%d\n", Config.Name, Config.Thresh)
}

func main2() {
	ctx := context.Background()
	client, err := etcdv3.New(
		etcdv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 3 * time.Second,
		},
	)
	if err != nil {
		panic(err)
	}

	InitGlobalConfig(ctx, client)
	go watch(ctx, client)

	for {
		UseConfig()
		time.Sleep(3 * time.Second)
	}
}

// go run .\micro_service\etcd\

/**
把etcd-v3.5.9-windows-amd64.zip解压后的目录放入环境变量PATH中，这样就可以在任意目录下使用etcd、etcdctl和etcdutl命令了

etcdctl命令示例：
etcdctl put dqq_name golang
etcdctl get dqq_name
etcdctl del dqq_name
*/
