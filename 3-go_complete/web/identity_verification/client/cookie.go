package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if resp, err := http.Get("http://localhost:5678/login"); err != nil {
		panic(err)
	} else {
		fmt.Println("response body")
		io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
		os.Stdout.WriteString("\n")
		loginCookies := resp.Cookies() //读取服务端返回的Cookie
		resp.Body.Close()
		if req, err := http.NewRequest(http.MethodGet, "http://localhost:5678/home", nil); err != nil {
			panic(err)
		} else {
			//下次请求再带上cookie
			for _, cookie := range loginCookies {
				fmt.Printf("receive cookie %s = %s\n", cookie.Name, cookie.Value)
				cookie.Value += "1" //修改cookie后认证不通过
				req.AddCookie(cookie)
			}
			client := &http.Client{Timeout: 1 * time.Second}
			if resp, err := client.Do(req); err != nil {
				fmt.Println(err)
			} else {
				defer resp.Body.Close()
				fmt.Println("response body")
				io.Copy(os.Stdout, resp.Body) //两个io数据流的拷贝
				os.Stdout.WriteString("\n")
			}
		}
	}
}
