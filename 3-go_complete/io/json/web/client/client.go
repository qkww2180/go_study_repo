package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	body := `{"Name":"大乔乔","age":18,"Birthday":"2023-09-30T08:44:50.5362744+08:00","CreatedAt":"2023-09-30"}`
	if resp, err := http.Post("http://127.0.0.1:5678", "application/json", strings.NewReader(body)); err != nil {
		log.Printf("http请求失败:%s", err)
	} else {
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
	}
}

// go run .\j_io\json\web\client\
