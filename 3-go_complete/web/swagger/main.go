package main

import (
	_ "dqq/web/swagger/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	engine := gin.Default()
	engine.GET("/swagger/*all", ginSwagger.WrapHandler(swaggerFiles.Handler)) //restful风格的url，一级目录是swagger，二级目录随意。必须打开这个路由才能访问http://127.0.0.1:5678/swagger/index.html。如果没的通过swag init命令生成对应的文档，系统也无法启动

	engine.GET("/get/:id", GetUser)         // get方法，restful风格的url
	engine.POST("/update_user", UpdateUser) // post方法

	engine.Run("127.0.0.1:5678")
}

// go run .\web\swagger\
// 浏览器打开 http://127.0.0.1:5678/swagger/index.html 可以查看文档，并在线调用接口
