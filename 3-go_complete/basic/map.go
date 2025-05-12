package main

import "fmt"

type Device struct {
	Name string
}

func main10() {
	m := map[string]Device{"1": Device{Name: "A"}}
	fmt.Println(m["2"])     //map的key不存在时返回零值，结构体的零值是空结构体，指针的零值是nil
	value, exists := m["2"] //map的key不存在时value是零值，exists是false
	fmt.Println(value, exists)
	// m["2"].Name = "B"
	// m["1"].Name = "B" //m["1"]可以读，但是不能写
	fmt.Println(value.Name) //空字符串
	value.Name = "B"        //value可以写

	for _, v := range m { //for range 取得的是拷贝
		v.Name = "B" //修改的是副本，不是map里的值
	}
	fmt.Println(m["1"].Name)

	for k := range m {
		value, _ := m[k] //go语言里所有的等号赋值都要发生拷贝
		value.Name = "B" //修改的是副本，不是map里的值
	}
	fmt.Println(m["1"].Name)

	mp := map[string]*Device{"1": &Device{Name: "A"}}
	for _, v := range mp { //取得的是指针的副本，但指针的副本和原指针指向同一块内存空间
		v.Name = "B" //修改了原始内存空间里的值
	}
	fmt.Println(mp["1"].Name)

	v, exists := mp["2"] //map的key不存在时返回零值，指针的零值是nil
	v.Name = "B"         //在nil上访问Name成员，会报 nil pointer dereference
}
