package json

import (
	"fmt"
	"time"
)

//自定义json里的时间格式，核心是自定义一个type，实现MarshalJSON和UnmarshalJSON这两个方法

var MyDateFormat = "2006-01-02"

type MyDate time.Time

func (d MyDate) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%s\"", time.Time(d).Format(MyDateFormat))
	return []byte(s), nil
}

// 要改变自己，必须传指针
func (d *MyDate) UnmarshalJSON(bs []byte) (err error) {
	now, err := time.ParseInLocation(`"`+MyDateFormat+`"`, string(bs), time.Local) //注意MyDateFormat前后还得加引号
	*d = MyDate(now)                                                               // 要改变自己
	return
}

// 在print(MyDate)时会调用String()方法
func (d MyDate) String() string {
	return time.Time(d).Format("2006-01-02")
}
