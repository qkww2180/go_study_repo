package main

import (
	"context"
	"dqq/micro_service/idl"
	"dqq/micro_service/util"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/unrolled/secure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// 没有单独启动grpc服务，只启动了一个http2服务
func main() {
	endpoint := "localhost:5679"
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(util.RootPath+"config/keys/server.crt", util.RootPath+"config/keys/server.key")
	if err != nil {
		panic(err)
	}

	grpcHandler := grpc.NewServer(grpc.Creds(creds))
	//绑定服务的实现
	idl.RegisterHelloHttpServer(grpcHandler, new(GreetService))

	// TLS认证
	creds2, err := credentials.NewClientTLSFromFile(util.RootPath+"config/keys/server.crt", "")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	// 连接grpc服务
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(creds2))
	if err != nil {
		panic(err)
	}

	gwHandler := runtime.NewServeMux() //实现了http.Handler接口
	//把post到/golang/hello上的请求转到grpc connection上
	if err := idl.RegisterHelloHttpHandler(context.Background(), gwHandler, conn); err != nil {
		panic(err)
	}

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") { //client走的是grpc
			log.Printf("走grpc")
			log.Printf("请求路径 %s", r.RequestURI)
			grpcHandler.ServeHTTP(w, r)
		} else { //client走的是http1.1
			log.Printf("走http1.1")
			gwHandler.ServeHTTP(w, r)
		}
	}))

	// 启动http服务
	secureMiddleware := secure.New(secure.Options{
		//把http://localhost:5678重定向到https://localhost:5678。这个选项其实可以不写，它是默认行为
		SSLRedirect: true,
		SSLHost:     endpoint,
	})
	app := secureMiddleware.Handler(router)
	//启动https（http+tls）服务
	http.ListenAndServeTLS(endpoint, util.RootPath+"config/keys/server.crt", util.RootPath+"config/keys/server.key", app)
}
