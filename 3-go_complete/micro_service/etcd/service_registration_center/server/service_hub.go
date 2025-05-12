package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"dqq/micro_service/etcd/service_registration_center"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

// 服务注册中心
type ServiceHub struct {
	client             *etcdv3.Client
	heartbeatFrequency int64 //server每隔几秒钟不动向中心上报一次心跳（其实就是续一次租约）
}

var (
	serviceHub *ServiceHub //该全局变量包外不可见，包外想使用时通过GetServiceHub()获得
	hubOnce    sync.Once   //单例模式需要用到一个once
)

// ServiceHub的构造函数，单例模式
func GetServiceHub(etcdServers []string, heartbeatFrequency int64) *ServiceHub {
	hubOnce.Do(func() {
		if serviceHub == nil {
			if client, err := etcdv3.New(
				etcdv3.Config{
					Endpoints:   etcdServers,
					DialTimeout: 3 * time.Second,
				},
			); err != nil {
				log.Fatalf("连接不上etcd服务器: %v", err) //发生log.Fatal时go进程会直接退出
			} else {
				serviceHub = &ServiceHub{
					client:             client,
					heartbeatFrequency: heartbeatFrequency, //租约的有效期
				}
			}
		}
	})
	return serviceHub
}

// 注册服务。 第一次注册向etcd写一个key，后续注册仅仅是在续约
//
// service 微服务的名称
//
// endpoint 微服务server的地址
func (hub *ServiceHub) Regist(service string, endpoint string, leaseID etcdv3.LeaseID) (etcdv3.LeaseID, error) {
	ctx := context.Background()
	if leaseID <= 0 {
		// 创建一个租约，有效期为heartbeatFrequency秒
		if lease, err := hub.client.Grant(ctx, hub.heartbeatFrequency); err != nil {
			log.Printf("创建租约失败：%v", err)
			return 0, err
		} else {
			key := strings.TrimRight(service_registration_center.SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
			// 服务注册
			if _, err = hub.client.Put(ctx, key, "", etcdv3.WithLease(lease.ID)); err != nil { //只需要key，不需要value
				log.Printf("写入服务%s对应的节点%s失败：%v", service, endpoint, err)
				return lease.ID, err
			} else {
				return lease.ID, nil
			}
		}
	} else {
		//续租
		if _, err := hub.client.KeepAliveOnce(ctx, leaseID); err == rpctypes.ErrLeaseNotFound { //续约一次，到期后还得再续约
			return hub.Regist(service, endpoint, 0) //找不到租约，走注册流程(把leaseID置为0)
		} else if err != nil {
			log.Printf("续约失败:%v", err)
			return 0, err
		} else {
			// log.Printf("服务%s对应的节点%s续约成功", service, endpoint)
			return leaseID, nil
		}
	}
}

// 注销服务
func (hub *ServiceHub) UnRegist(service string, endpoint string) error {
	ctx := context.Background()
	key := strings.TrimRight(service_registration_center.SERVICE_ROOT_PATH, "/") + "/" + service + "/" + endpoint
	if _, err := hub.client.Delete(ctx, key); err != nil {
		log.Printf("注销服务%s对应的节点%s失败: %v", service, endpoint, err)
		return err
	} else {
		log.Printf("注销服务%s对应的节点%s", service, endpoint)
		return nil
	}
}
