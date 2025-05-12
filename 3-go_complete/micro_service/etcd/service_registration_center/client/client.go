package main

import (
	"context"
	"dqq/micro_service/etcd/service_registration_center"
	"dqq/micro_service/idl"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverConn = sync.Map{}

func getClient() idl.HelloServiceClient {
	hub := GetServiceHub(service_registration_center.ETCD_CLUSTER)
	//根据service name获取所有可用的servers
	servers := hub.GetServiceEndpointsWithCache(service_registration_center.HELLO_SERVICE)
	if len(servers) == 0 {
		log.Printf("服务%s没有对应的实例", service_registration_center.HELLO_SERVICE)
		return nil
	} else {
		//采用某种负载均衡算法，选择一个server
		idx := rand.Intn(len(servers)) //这里采用随机法
		server := servers[idx]
		fmt.Println(server)
		if client, exists := serverConn.Load(server); exists {
			return client.(idl.HelloServiceClient)
		} else {
			// 连接到GRPC服务端
			conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Printf("连接GRPC服务端失败 %v\n", err)
				return nil
			}
			// log.Printf("连接到server %s\n", server)
			client := idl.NewHelloServiceClient(conn)
			serverConn.Store(server, client)
			return client
		}
	}
}

func rpc() {
	client := getClient()
	if client == nil {
		log.Printf("RPC调用失败，无法连接到server")
	}
	// log.Println("获得client")
	client.Login(context.Background(), &idl.LoginRequest{})
	// log.Println("server返回结果")
}

func main() {
	for {
		rpc()
		time.Sleep(time.Second)
	}
}

// go run .\micro_service\etcd\service_registration_center\client
