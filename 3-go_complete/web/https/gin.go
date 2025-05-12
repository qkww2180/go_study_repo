package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure" //需要用到这个第三方库
)

/**
https = http + tls
首先，生成证书和key。
go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
go run "D:/Program Files/Go/src/crypto/tls/generate_cert.go" --host="localhost"
在当前目录下会生成cert.pem和key.pem这两个文件。我把这两个文件移到了config/keys目录下(这个操作不是必须的)

在《微服务》课程里面使用的是以下步骤，更清晰。
1. 生成server的私钥
openssl genrsa -out config/keys/server.priv 2048
2. 生成证书
openssl req -x509 -new -nodes -key config/keys/server.priv -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost" -days 3650 -out config/keys/server.crt

需要把自己造的证书导入到操作系统受信任的证书列表里，否则浏览器会拦截。操作方法参见《微服务》课程
*/

func main() {
	secureMiddleware := secure.New(secure.Options{
		//把http://localhost:5678重定向到https://localhost:5678。这个选项其实可以不写，它是默认行为
		SSLRedirect: true,
		SSLHost:     "localhost:5678",
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)
			if err != nil {
				c.Abort()
				return
			}
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	engine := gin.Default()
	engine.Use(secureFunc)

	engine.GET("/", func(c *gin.Context) {
		c.String(200, "欢迎来到HTTP Secure的世界")
	})
	//启动https（http+tls）服务
	// engine.RunTLS("localhost:5678", "config/keys/cert.pem", "config/keys/key.pem")
	engine.RunTLS("localhost:5678", "config/keys/server.crt", "config/keys/server.priv")
}

// go run .\web\https\
