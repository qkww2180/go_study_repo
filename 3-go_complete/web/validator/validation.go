package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10" //注意要用新版本v10
)

type Student struct {
	Name  string `form:"name" binding:"required"` //required:必须上传name参数。form可以绑定formdata和url里的参数
	Score int    `form:"score" binding:"gt=0"`    //score必须为正数

	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`       //自定义验证before_today，日期格式东8区
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"` //毕业时间要晚于入学时间
}

// 自定义验证器
var beforeToday validator.Func = func(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok { //通过反射获得结构体Field的值
		today := time.Now()
		if date.Before(today) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func processErr(err error) string {
	if err == nil {
		return ""
	}

	//给Validate.Struct()函数传了一个非法的参数
	invalid, ok := err.(*validator.InvalidValidationError)
	if ok {
		return fmt.Sprintf("param error: %v", invalid)
	}

	//ValidationErrors是一个错误切片，它保存了每个字段违反的每个约束信息
	validationErrs := err.(validator.ValidationErrors)
	msgs := make([]string, 0, 3)
	for _, validationErr := range validationErrs {
		msgs = append(msgs, fmt.Sprintf("field %s 不满足条件 %s", validationErr.Field(), validationErr.Tag()))
	}
	return strings.Join(msgs, ";")
}

func main() {
	engine := gin.Default()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	engine.GET("/", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBind(&stu); err != nil { //在绑定参数的同时，完成合法性校验
			msg := processErr(err)
			ctx.String(http.StatusBadRequest, "parse parameter failed. "+msg) //校验不符合时，返回哪里不符合
		} else {
			ctx.JSON(http.StatusOK, stu) //校验通过时，返回一个json
		}
	}) //http://localhost:5678?name=zcy&score=1&enrollment=2021-08-23&graduation=2021-09-23

	engine.Run("127.0.0.1:5678")
}
