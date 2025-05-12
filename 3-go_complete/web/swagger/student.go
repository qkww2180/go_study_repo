package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`   //用户ID
	Name string `json:"name"` //姓名
	Age  int    `json:"age"`  //年龄
}

// 参数类型有：header, body(正常的post参数), query(正常的get参数), path(restful风格的参数)

//	@Summary	获取用户信息
//	@Produce	json
//	@Param		id	path		int		true	"用户ID"
//	@Success	200	{object}	User	"成功"
//	@Failure	400	{object}	string	"参数错误"
//	@Failure	500	{object}	string	"内部错误"
//	@Router		/get/{id} [GET]
func GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id") //从Restful风格的url中获取参数
	if id, err := strconv.Atoi(idStr); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
	} else {
		// ctx.String(http.StatusInternalServerError, "内部错误")//比如读数据库失败
		ctx.JSON(http.StatusOK, User{ID: id, Name: "大乔乔", Age: 18})
	}
}

//	@Summary	更新用户信息
//	@Produce	json
//	@Param		uer	body		User	true	"用户信息"
//	@Success	200	{object}	string	"更新成功"
//	@Failure	400	{object}	string	"参数错误"
//	@Failure	500	{object}	string	"内部错误"
//	@Router		/update_user [POST]
func UpdateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
	} else {
		// ctx.String(http.StatusInternalServerError, "内部错误")//比如写数据库失败
		ctx.String(http.StatusOK, "更新成功")
	}
}

// swag fmt .\web\swagger\    	格式化swag注释
// cd .\web\swagger\ ; swag init   在.\web\swagger目录下会生成
// ./docs
//   |---docs.go				注意这里有一个go文件,go文件里有个init()函数
//   |---swagger.json
//   |---swagger.yaml
