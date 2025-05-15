package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"dqq/micro_service/etcd/service_registration_center"

	etcdv3 "go.etcd.j_io/etcd/client/v3"
)

// 服务注册中心
type ServiceHub struct {
	client        *etcdv3.Client
	endpointCache sync.Map //维护每一个service下的所有servers
	watched       sync.Map
}

var (
	serviceHub *ServiceHub //该全局变量包外不可见，包外想使用时通过GetServiceHub()获得
	hubOnce    sync.Once   //单例模式需要用到一个once
)

// ServiceHub的构造函数，单例模式
func GetServiceHub(etcdServers []string) *ServiceHub {
	hubOnce.Do(func() {
		if serviceHub == nil {
			if client, err := etcdv3.New(
				etcdv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				},
			); err != nil {
				log.Fatalf("连接不上etcd服务器: %v", err)
			} else {
				serviceHub = &ServiceHub{
					client:        client,
					endpointCache: sync.Map{},
					watched:       sync.Map{},
				}
			}
		}
	})
	return serviceHub
}

// 服务发现。client每次进行RPC调用之前都查询etcd，获取server集合，然后采用负载均衡算法选择一台server。或者也可以把负载均衡的功能放到注册中心，即放到getServiceEndpoints函数里，让它只返回一个server
func (hub *ServiceHub) getServiceEndpoints(service string) []string {
	ctx := context.Background()
	prefix := strings.TrimRight(service_registration_center.SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	if resp, err := hub.client.Get(ctx, prefix, etcdv3.WithPrefix()); err != nil { //按前缀获取key-value
		log.Printf("获取服务%s的节点失败: %v", service, err)
		return nil
	} else {
		endpoints := make([]string, 0, len(resp.Kvs))
		for _, kv := range resp.Kvs {
			path := strings.Split(string(kv.Key), "/") //只需要key，不需要value
			// fmt.Println(string(kv.Key), path[len(path)-1])
			endpoints = append(endpoints, path[len(path)-1])
		}
		log.Printf("刷新%s服务对应的server -- %v\n", service, endpoints)
		return endpoints
	}
}

func (hub *ServiceHub) watchEndpointsOfService(service string) {
	if _, exists := hub.watched.LoadOrStore(service, true); exists {
		return //监听过了，不用重复监听
	}
	ctx := context.Background()
	prefix := strings.TrimRight(service_registration_center.SERVICE_ROOT_PATH, "/") + "/" + service + "/"
	ch := hub.client.Watch(ctx, prefix, etcdv3.WithPrefix()) //根据前缀监听，每一个修改都会放入管道ch
	log.Printf("监听服务%s的节点变化", service)
	go func() {
		for response := range ch { //遍历管道。这是个死循环，除非关闭管道
			for _, event := range response.Events { //每次从ch里取出来的是事件的集合
				path := strings.Split(string(event.Kv.Key), "/")
				if len(path) > 2 {
					service := path[len(path)-2]
					// 跟etcd进行一次全量同步
					endpoints := hub.getServiceEndpoints(service)
					if len(endpoints) > 0 {
						hub.endpointCache.Store(service, endpoints) //查询etcd的结果放入本地缓存
					} else {
						hub.endpointCache.Delete(service) //该service下已经没有endpoint
					}
				}
			}
		}
	}()
}

// 服务发现。把第一次查询etcd的结果缓存起来，然后安装一个Watcher，仅etcd数据变化时更新本地缓存。这样可以降低etcd的访问压力
func (hub *ServiceHub) GetServiceEndpointsWithCache(service string) []string {
	hub.watchEndpointsOfService(service) //监听etcd的数据变化，及时更新本地缓存
	if endpoints, exists := hub.endpointCache.Load(service); exists {
		return endpoints.([]string)
	} else {
		endpoints := hub.getServiceEndpoints(service)
		if len(endpoints) > 0 {
			hub.endpointCache.Store(service, endpoints) //查询etcd的结果放入本地缓存
		}
		return endpoints
	}
}
