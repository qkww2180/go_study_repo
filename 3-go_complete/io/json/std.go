package json

import (
	"time"
)

// 序列化就是把变量（包括基础类型的变量和复杂的结构体变量）转为二进制流（二进制流可以跟[]byte等价，[]byte又可以跟string等价），以便写入磁盘或发送到网络上。反序列化跟该过程相反。

type User struct {
	Name      string    //默认的json字段名跟原始变量名保持一致
	Age       int       `json:"age"`
	height    float32   //不可导出成员不会被序列化（该变量的值不会被导出到磁盘或网络上），否则就违背了“不可导出”的本意
	Birthday  time.Time //格式： 2023-09-29T20:14:11.7074482+08:00
	CreatedAt MyDate    //格式： 2023-09-29
}

// 标准库json序列化背后使用的核心技术是反射。通过反射可以在运行时动态获得结构体成员变量的名称、(json)tag、是否可导出，可以获取成员变量的值，还可以调用结构体的方法(比如MarshalJSON和UnmarshalJSON)
