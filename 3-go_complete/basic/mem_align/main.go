package main

import (
	"fmt"
	"reflect"
)

/*
struct的内存对齐问题
*/
func main() {
	type A struct {
		Sex      bool   //offset  0
		Height   uint16 //offset  2   它本身占2个字节，它只从偶数位开始，所以即使sex只占1个字节，height的offset也是2
		Addr     byte   //offset  4
		Position byte   //offset  5
		Age      int64  //offset  8   它本身占8个字节，它只从8的整倍数开始，即使之前的字段只占6个字节
		Weight   uint16 //offset  16
	} //整个结构体占24个字节，因为结构体的内存对齐原则是：最大成员变量的整倍数。Age是最大的成员变量，所以结构体的内存占用是8B的整倍数
	t := reflect.TypeOf(A{})
	fmt.Println("结构体内存开销", t.Size()) //24B

	fieldNum := t.NumField() //成员变量的个数，包括未导出成员
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		fmt.Printf("%s offset %d\n",
			field.Name,   //变量名称
			field.Offset, //相对于结构体首地址的内存偏移量，string类型会占据16个字节
		)
	}
}
