package main

import (
	"fmt"
	"math"
)

func main2() {
	var MyName int      //变量声明
	fmt.Print(MyName)   //使用变量
	fmt.Println(MyName) //使用变量
	var a int = 8
	var b = a //自动推断b的类型
	_ = b
	c := b // 第一次出现 : 声明  var
	a = c  //第二次出现不能加:和var

	var (
		d uint16
		e int8
		f float32
		g float64
	)
	a = -5
	d = 05    //前缀0表示八进制
	a = 0o57  //前缀0O表示八进制
	a = 0xab3 //前缀0x表示十六进制
	a = 5_0_123_7
	a = 13_000_000 //13M
	f = 1.4324325346543654
	g = 34.
	m := 34.          //float64
	var n bool = true //默认值为false
	_, _, _, _ = d, e, f, m

	fmt.Printf("a=%d, g=%.2f, n=%t\n", a, g, n) //2位小数
	fmt.Printf("f=%g, f=%e\n", f, f)
	fmt.Printf("f=%[1]f, g=%[2]f, g=%[2]g, f=%[1]e\n", f, g)

	cast()
}

// 强制类型转换
func cast() {
	//高精度向低精度转换，数字很小时这种转换没问题
	var ua uint64 = 1
	i8 := int8(ua)
	fmt.Printf("i8=%d\n", i8)

	//最高位的1变成了符号位
	ua = uint64(math.MaxUint64)
	i64 := int64(ua)
	fmt.Printf("i64=%d\n", i64) //-1。负数用补码形式存储   00000000000001  1111111111111

	//位数丢失
	ui32 := uint32(ua)
	fmt.Printf("ui32=%d\n", ui32)

	//单个字符可以转为int
	var i int = int('中') //'中'是字符，“中”是字符串
	fmt.Printf("i=%d\n", i)

	//bool和int不能相互转换

	//byte和int可以互相转换
	var by byte = byte(i)
	i = int(by)
	fmt.Printf("i=%d\n", i)

	//float和int可以互相转换，小数位会丢失
	var ft float32 = float32(i)
	i = int(ft)
	fmt.Printf("i=%d\n", i)
}
