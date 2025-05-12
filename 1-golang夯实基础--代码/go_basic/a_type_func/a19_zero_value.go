package main

import (
	"fmt"
)

// zero value，零值，默认值
func main19() {
	var i int //各种整型(包括rune)，0
	fmt.Printf("i=%v\n", i)
	var f float32 //各种浮点型，0
	fmt.Printf("f=%v\n", f)
	var b byte //0
	fmt.Printf("b=%v\n", b)
	var bl bool //false
	fmt.Printf("bl=%v\n", bl)
	var s string // ""
	fmt.Printf("s=%v\n", s)
	var p *int //指针，nil
	fmt.Printf("p=%v\n", p)
	type User struct {
		Gender bool
		Age    int
	}
	var u User //结构体的每一个成员变量分别为相对数据类型的“0值”，可能会存在嵌套
	fmt.Printf("u=%v\n", u)
	var err error //接口，nil
	fmt.Printf("err=%v\n", err)
	err = nil
	var arr [3]int //数组，所有成员全是0值
	fmt.Printf("arr=%v\n", arr)
	var slc []string            //引用类型，nil
	fmt.Printf("slc=%v\n", slc) //长度为0、容量为0的空切片
	var mp map[int]bool         //引用类型，nil
	fmt.Printf("mp=%v\n", mp)   //长度为0的空map
	// mp[2] = true    //panic: assignment to entry in nil map。map必须经过初始化才能进行读写
	var ch chan int //引用类型，nil
	fmt.Printf("ch=%v\n", ch)
	// var fk func(bool) error //函数类型，nil
	// fmt.Printf("fk=%v\n", fk)

	// 空接口，类型断言失败，返回零值
	var ifc any
	v1, ok := ifc.(float64) // nil接口向任何类型断言都会失败，不会发生空指针error
	fmt.Println(v1, ok)     // 0 false
	ifc = 3.14              // go语言里的字面量小数是float64，不是float32
	v2, ok := ifc.(float32) // 0 false
	fmt.Println(v2, ok)
	v3, ok := ifc.(float64) // 3.14 true
	fmt.Println(v3, ok)

	mapValue, exists := mp[7] //key不存在时，读出的value是0值。通过exists来告诉你是不是0值
	fmt.Printf("mapValue=%v\n", mapValue)
	_ = exists

	ch = make(chan int, 10)
	close(ch)
	chValue, ok := <-ch //channel为空，且channel被关闭后，读出的value是0值。通过ok来告诉你是不是0值
	fmt.Printf("chValue=%v\n", chValue)
	_ = ok

	// err:=gorm.Select("*").Where("age=0").First(&u).Error() 查不到结果时u也是0值。通过err来告诉你是不是0值
}
