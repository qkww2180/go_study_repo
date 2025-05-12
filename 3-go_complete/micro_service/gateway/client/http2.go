package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	endpoint = "localhost:5679" // 证书里Common Name指定的域名是localhost，所以这里只能使用localhost
)

func grpcClient2() {
	// 连接到GRPC服务端
	dialCred, err := credentials.NewClientTLSFromFile(util.RootPath+"config/keys/server.crt", "") //网关在请求grpc服务时需要使用server的证书
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(dialCred))
	if err != nil {
		fmt.Printf("连接GRPC服务端失败 %v\n", err)
		return
	}
	defer conn.Close()
	client := idl.NewHelloHttpClient(conn)

	// 执行RPC调用并打印收到的响应数据，指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := client.Greeting(ctx, request)
	if err != nil {
		fmt.Printf("远程调用失败 %v\n", err)
		return
	} else {
		fmt.Println(resp.Message)
	}
}

func httpClient2() {
	//Server的证书(这里用的是自签名证书)
	serverCrt, err := tls.LoadX509KeyPair(util.RootPath+"config/keys/server.crt", util.RootPath+"config/keys/server.key")
	if err != nil {
		panic(err)
	}

	//CA的证书(这里Server自己充当CA)
	caCrt, _ := os.ReadFile(util.RootPath + "config/keys/server.crt")
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt) //把CA的证书放入池中

	// 构造http.Client时传入证书
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{ //因为要发起https请求，所以需要指定TLS配置
				RootCAs:      pool,
				Certificates: []tls.Certificate{serverCrt},
			},
		},
	}

	reqStr, err := sonic.MarshalString(request)
	if err != nil {
		return
	}
	reader := strings.NewReader(reqStr)

	// 在.proto文件里指定的请求方法是post，路径是"/golang/hello"
	if resp, err := client.Post("https://"+endpoint+"/golang/hello", "application/json", reader); err != nil { //请求是json格式
		panic(err)
	} else {
		defer client.CloseIdleConnections()
		defer resp.Body.Close() //注意：一定要调用resp.Body.Close()，否则会协程泄漏（同时引发内存泄漏）
		bs, _ := io.ReadAll(resp.Body)
		var response idl.GreetResopnse
		sonic.Unmarshal(bs, &response)
		fmt.Println(response.Message)
	}
}

func main() {
	grpcClient2()
	httpClient2()
}

// go run .\micro_service\gateway\client\
