package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := http.Client{
		Timeout: 1 * time.Second, //小于1秒，导致请求超时，会触发Server端的http.Request.Context的Done
	}
	if resp, err := client.Get("http://127.0.0.1:5678/"); err == nil {
		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)
		if bs, err := io.ReadAll(resp.Body); err == nil {
			fmt.Println(string(bs))
		}
	} else {
		fmt.Println(err) //Get "http://127.0.0.1:5678/": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
	}
}

// go run ./basic/timeout/client
