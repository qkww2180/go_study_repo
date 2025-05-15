package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Student struct {
	Name   string
	Age    int
	Height float32
}

type Request struct {
	StudentId string `json:"student_id"`
}

// 从redis上根据studentId获取Student实体
func GetStudentInfo(studentId string) Student {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.TODO()
	stu := Student{}
	for field, value := range client.HGetAll(ctx, "学生1").Val() {
		if field == "Name" {
			stu.Name = value
		} else if field == "Age" {
			age, err := strconv.Atoi(value)
			if err == nil {
				stu.Age = age
			}
		} else if field == "Height" {
			height, err := strconv.ParseFloat(value, 10)
			if err == nil {
				stu.Height = float32(height)
			}
		}
	}
	return stu
}

func GetName(ctx *gin.Context) {
	param := ctx.Query("student_id") //从Get请求中获取参数
	if len(param) == 0 {             //没有指定student_id参数
		ctx.String(http.StatusBadRequest, "please indidate student_id") //StatusBadRequest即400
		return
	}
	stu := GetStudentInfo(param)
	ctx.String(http.StatusOK, stu.Name) //StatusOK即200
	// ctx.JSON(h_http.StatusOK,stu)
	return
}

func GetAge(ctx *gin.Context) {
	param := ctx.PostForm("student_id") //从post form中获取参数
	if len(param) == 0 {                //没有指定student_id参数
		ctx.String(http.StatusBadRequest, "please indidate student_id") //StatusBadRequest即400
		return
	}
	stu := GetStudentInfo(param)
	ctx.String(http.StatusOK, strconv.Itoa(stu.Age)) //StatusOK即200
	// ctx.JSON(h_http.StatusOK,stu)
	return
}

func GetHeight(ctx *gin.Context) {
	var request Request
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.String(http.StatusBadRequest, "please indidate student_id in json")
		return
	}
	stu := GetStudentInfo(request.StudentId)
	ctx.String(http.StatusOK, strconv.FormatFloat(float64(stu.Height), 'f', 1, 64)) //保留1位小数
	// ctx.JSON(h_http.StatusOK, stu)
	return
}

func main() {
	engine := gin.Default()
	engine.GET("/get_name", GetName)
	engine.POST("/get_age", GetAge)
	engine.POST("/get_height", GetHeight)

	err := engine.Run("127.0.0.1:2345")
	if err != nil {
		panic(err)
	}
}

// go run ./gin
