package main

import "fmt"

const PI = 100 //全局变量

func ellop() {
	fmt.Println(PI) //使用的是全局常量PI
}

func main4() {
	var L int = 100
	_ = L
	// {}限定了变量的作用域
	{
		const PI float32 = 3.14   //常量声明时必须赋值，以后不能再更改
		fmt.Printf("PI=%f\n", PI) //就近原则，使用的是局部常量PI
		var L int = 20
		fmt.Println(L) //就近原则
	}
	{
		const (
			PI = 3.14
			E  = 2.71
		)
		fmt.Printf("PI=%f, E=%f\n", PI, E)
	}
	{
		const (
			a = 100
			b //100，跟上一行的值相同
			c //100，跟上一行的值相同
			d = 9
			e //9，跟上一行的值相同
		)
		fmt.Printf("a=%d, b=%d, c=%d, d=%d, e=%d\n", a, b, c, d, e)
	}
	{
		const (
			a = iota //0
			b        //1
			c        //2
			d        //3
		)
		fmt.Printf("a=%d, b=%d, c=%d, d=%d\n", a, b, c, d)
	}
	{
		const (
			a = iota //0
			b        //1
			_        //2
			d        //3
		)
		fmt.Printf("a=%d, b=%d, d=%d\n", a, b, d)
	}
	{
		const (
			a = iota //0
			b = 50   //50
			c = iota //2
			d        //3
		)
		fmt.Printf("a=%d, b=%d, c=%d, d=%d\n", a, b, c, d)
	}
	{
		const (
			_  = iota             // iota =0
			KB = 1 << (10 * iota) // iota =1
			MB = 1 << (10 * iota) // iota =2
			GB = 1 << (10 * iota) // iota =3
			TB = 1 << (10 * iota) // iota =4
		)
		fmt.Printf("KB=%d, MB=%d, GB=%d, TB=%d\n", KB, MB, GB, TB)
	}
	{
		const (
			a, b = iota + 1, iota + 2 //1,2  iota=0
			c, d                      //2,3  iota=1
			e, f                      //3,4  iota=2
		)
		fmt.Printf("a=%d, b=%d, c=%d, d=%d, e=%d, f=%d\n", a, b, c, d, e, f)
	}
	ellop()
}
