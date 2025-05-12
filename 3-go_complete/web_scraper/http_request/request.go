package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var serverAddress = "http://127.0.0.1:5678"

// 不带参数的get请求
func getWithoutParam() {
	if resp, err := http.Get(serverAddress); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		/**
		具体看一下http协议
		*/
		fmt.Printf("response proto: %s\n", resp.Proto)
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response header")
		for key, values := range resp.Header {
			fmt.Printf("%s: %v\n", key, values)
		}
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n")
	}
}

// 带参数的get请求
func getWithoParam() {
	if resp, err := http.Get(serverAddress + "?a=4&b=6"); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		/**
		具体看一下http协议
		*/
		fmt.Printf("response proto: %s\n", resp.Proto)
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response header")
		for key, values := range resp.Header {
			fmt.Printf("%s: %v\n", key, values)
		}
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n")
	}
}

// HEAD类似于GET，但HEAD方法只能取到http response报文头部，取不到resp.Body。head请求通过用于验证一个url是否存活
func head() {
	if resp, err := http.Head(serverAddress); err != nil {
		log.Printf("http请求失败:%s", err) //如果url不存在，此处会报错
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		/**
		具体看一下http协议
		*/
		fmt.Printf("response proto: %s\n", resp.Proto)
		fmt.Printf("response status: %s\n", resp.Status) //状态码为200就说明url存活
		fmt.Println("response header")
		for key, values := range resp.Header {
			fmt.Printf("%s: %v\n", key, values)
		}
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //resp.Body为空
		os.Stdout.WriteString("\n")
	}
}

// body为空的http请求
func postWithoutBody() {
	if resp, err := http.Post(serverAddress+"?a=4&b=6", "text/plain", nil); err != nil { //请求body为空，参数放在url里，当然url里也可以不放参数(如果server端不需要传的话)
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n")
	}
}

// post请求，request body是个普通的字符串
func postStringBody() {
	body := "hello"
	if resp, err := http.Post(serverAddress, "text/plain", strings.NewReader(body)); err != nil { //Content-Type为text/plain，表示一个朴素的字符串
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		/**
		具体看一下http协议
		*/
		fmt.Printf("response proto: %s\n", resp.Proto)
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response header")
		for key, values := range resp.Header {
			fmt.Printf("%s: %v\n", key, values)
		}
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n")
	}
}

type Request struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}
type Response struct {
	Age   int `json:"age"`
	Score int `json:"score"`
}

// post请求，request body和request body都是json，json本质上也是字符串，所以跟postStringBody()相同，只是前后加了json序列化和反序列化
func postJsonBody() {
	request := &Request{Id: 3, Token: "fj39834"}
	body, _ := json.Marshal(request)                                                                  //body是byte切片
	if resp, err := http.Post(serverAddress, "application/json", bytes.NewReader(body)); err != nil { //Content-Type为text/json，表示一个json字符串
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		bs, _ := io.ReadAll(resp.Body) //把body流全部读到byte切片里
		var response Response
		json.Unmarshal(bs, &response)
	}
}

// post请求，请求参数放在form表单里
func postForm() {
	if resp, err := http.PostForm(serverAddress, url.Values{"name": []string{"golang"}, "age": []string{"18"}}); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n")
	}
}

func postForm2() {
	request, _ := http.NewRequest(http.MethodPost, serverAddress, nil)
	form := url.Values{}
	form.Add("name", "golang")
	form.Add("age", "18")
	request.Form = form
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	//发送http请求
	if resp, err := client.Do(request); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n")
	}
}

// 通过http.Client发送请求是一种万能的方式，可以涵盖以上的所有方式，并且可以设计Header和Cookie(以上方法不能设置)
func httpClient() {
	request := &Request{Id: 3, Token: "fj39834"}
	body, _ := json.Marshal(request)
	reader := bytes.NewReader(body)
	//构造http request
	req, _ := http.NewRequest(http.MethodPost, serverAddress, reader)
	//向http request里添加header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
	//向http request里添加cookie
	req.AddCookie(&http.Cookie{
		Name:   "auth",
		Value:  "pass",
		Path:   "/",
		Domain: "localhost",
	})
	//构造http client
	client := &http.Client{
		Timeout: 500 * time.Millisecond, //设置请求超时
		Transport: &http.Transport{ //设置跟网络传输相关的参数
			Proxy: http.ProxyFromEnvironment, //从环境变量中读取代理，可以是hhtp、https、socks5
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	//发送http请求
	if resp, err := client.Do(req); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n")
	}
}
