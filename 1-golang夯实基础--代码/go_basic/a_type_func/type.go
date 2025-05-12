package main

import "fmt"

type Age int    //Age和int可以互相做强制类型转换。存储的数据类型（成员变量）是一样的，但行为（成员方法）是不一样
type Tall = int //Tall和int完全等价，不需要显式做类型转换

type Odt Shiper // Shiper实现了Steamer接口

// Age实现了Transporter接口
func (Age) move(src string, dest string) (int, error) {
	return 0, nil
}

func (Age) whistle(n int) int {
	return 0
}

func main21() {
	var a int
	var b Age
	var c Tall

	transport("BJ", "SH", b)

	c = a + Tall(10)
	b = Age(a) + 10 //字面量不需要显式转为Age类型

	fmt.Println(a, b, c)

	var o Odt
	o.tonage = 100
	o.name = "dqq"
	seaTransport("BJ", "SH", Shiper(o)) //Odt没有实现Steamer接口
}
