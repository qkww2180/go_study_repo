package orm

import (
	"fmt"
	"reflect"
	"time"
)

type Account struct {
	Name string
}

type Student struct {
	Id         int
	Province   string `gorm:"column:province" json:"province"` //可以包含多个Tag，用空格分开
	City       string
	Address    string
	Account    //匿名成员，这样相当于Student“继承”了Account
	Enrollment time.Time
	Gender     int `gorm:"-"`
}

func (Student) Hello(name string) string {
	return "hello " + name
}

// 打印结构体的成员变量信息
func PrintFieldInfo(object any) {
	tp := reflect.TypeOf(object) //通过reflect.Type获取类型相关的信息
	fieldNum := tp.NumField()    //成员变量的个数，包括未导出成员
	for i := 0; i < fieldNum; i++ {
		field := tp.Field(i)
		fmt.Printf("%d %s offset %d anonymous %t type %s exported %t gorm tag %s json tag %s\n", i,
			field.Name,            //变量名称
			field.Offset,          //相对于结构体首地址的内存偏移量，string类型会占据16个字节
			field.Anonymous,       //是否为匿名成员
			field.Type,            //数据类型，reflect.Type类型
			field.IsExported(),    //包外是否可见（即是否以大写字母开头）
			field.Tag.Get("gorm"), //获取成员变量后面``里面定义的gorm
			field.Tag.Get("json"), //获取成员变量后面``里面定义的tag
		)
	}
	fmt.Println()
}

// 打印结构体的成员变量的值，然后修改它
func PrintFieldValue(object any) {
	tp := reflect.TypeOf(object)                    //通过reflect.Type获取类型相关的信息
	vl := reflect.ValueOf(object)                   //通过reflect.Value获取、修改原始数据类型里的值
	if _, ok := tp.Elem().FieldByName("City"); ok { //先确保存在City这个Field
		city := vl.Elem().FieldByName("City") //通过Elem()把vl转为非指针形式
		fmt.Println("修改前：" + city.String())   //如果City不是string类型，运行时才会报错
		if city.CanAddr() && city.CanSet() {  //由于传进来的是指针，所以可寻址；由于City是大写开头可导出的，所以是CanSet
			city.SetString("北京") //此处使用SetFloat()也不会编译出错，到运行时才能发现错误。这也是使用反射难以测试，容易出线上问题的原因
		}
	}
}

// 执行结构体的方法
func ExecMethod(object any) {
	vl := reflect.ValueOf(object)      //通过reflect.Value获取、修改原始数据类型里的值
	method := vl.MethodByName("Hello") //MethodByName()通过Name返回类的成员变量
	arg := reflect.ValueOf("大乔乔")
	results := method.Call([]reflect.Value{arg}) //方法的返回值可能有多个
	result := results[0].Interface().(string)    //取得第一个返回值，断言它是string类型
	fmt.Println(result)
}
