package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket" //该项目带了一个聊天室的demo https://github.com/gorilla/websocket/tree/master/examples/chat
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 1 * time.Second,
	ReadBufferSize:   100,
	WriteBufferSize:  100,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	frontName []byte //前端的名字，用于展示在消息前面
}

// 从websocket连接里读出数据，发给hub
func (client *Client) read() {
	defer func() { //收尾工作
		client.hub.unregister <- client //从hub那注销client
		fmt.Printf("%s offline\n", client.frontName)
		fmt.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //关闭websocket管道
	}()
	for {
		_, message, err := client.conn.ReadMessage() //如果前端主动断开连接，该行会报错，for循环会退出。注销client时，hub那儿会关闭client.send管道
		if err != nil {
			break //只要ReadMessage失败，就关闭websocket管道、注销client，退出
		} else {
			//换行符用空格替代，bytes.TrimSpace把首尾连续的空格去掉
			message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
			if len(client.frontName) == 0 {
				client.frontName = message //约定：从浏览器读到的第一条消息代表前端的身份标识，该信息不进行广播
				client.hub.broadcast <- []byte(fmt.Sprintf("%s online\n", string(client.frontName)))
			} else {
				//要广播的内容前面加上front的名字
				client.hub.broadcast <- bytes.Join([][]byte{client.frontName, message}, []byte(": ")) //从websocket连接里读出数据，发给hub的broadcast
			}
		}
	}
}

// 从hub的broadcast那儿读限数据，写到websocket连接里面去
func (client *Client) write() {
	defer func() {
		fmt.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //给前端写数据失败，就可以关系连接了
	}()

	for {
		msg, ok := <-client.send
		if !ok {
			fmt.Println("管道已经被关闭")
			client.conn.WriteMessage(websocket.CloseMessage, []byte("bye bye"))
			return
		} else {
			err := client.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("向浏览器发送数据失败:%v\n", err)
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) //http升级为websocket协议
	if err != nil {
		fmt.Printf("upgrade error: %v\n", err)
		return
	}
	fmt.Printf("connect to client %s\n", conn.RemoteAddr().String())
	//每来一个前端请求，就会创建一个client
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//分别设置读写超时
	// client.conn.SetWriteDeadline(time.Now().Add(24 * time.Hour))
	// client.conn.SetReadDeadline(time.Now().Add(24 * time.Hour))
	//向hub注册client
	client.hub.register <- client

	//启动子协程，运行ServeWs的协程退出后子协程也不会能出
	//websocket是全双工模式，可以同时read和write。即websocket的connection支持并发操作
	go client.read()
	go client.write()
}
