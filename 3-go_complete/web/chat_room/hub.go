package main

type Hub struct {
	broadcast  chan []byte          //需要广播的消息
	clients    map[*Client]struct{} //维护所有Client
	register   chan *Client         //Client注册请求通过管道来接收
	unregister chan *Client         //Client注销请求通过管道来接收
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]struct{}),
		//以下都是无缓冲管道，确保client和hub的交互是同步的，即client的请求hub处理完之后client才执行下一步工作
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (hub *Hub) Run() {
	//通过select机制避免了并发修改了clients；避免了遍历clients的同时修改clients。
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = struct{}{} //注册client
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok { //防止重复注销
				delete(hub.clients, client) //注销client
				close(client.send)          //hub从此以后不需要再向该client广播消息了
			}
		case msg := <-hub.broadcast:
			for client := range hub.clients {
				client.send <- msg
			}
		}
	}
}
