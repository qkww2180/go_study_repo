package main

import (
	"dqq/io/json"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/", func(ctx *gin.Context) {
		var user json.User
		if err := ctx.ShouldBindJSON(&user); err != nil { //BindJSON背后还是调用了标准库的json.Unmarshal
			log.Printf("解析参数失败: %s", err)
		} else {
			ctx.JSON(200, user)
		}
	})
	router.Run("127.0.0.1:5678")
}

// go run .\io\json\web\gin_server\
