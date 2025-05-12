package main

import "fmt"

// 接口是一组行为规范的集合
type EmptyInterface interface{} //EmptyInterface的go语言标准库里的interface{}、any本质是一样，但是不等价（不能互换）

func SumI(args ...EmptyInterface) int {
	rect := 0
	for _, ele := range args {
		switch v := ele.(type) { //.(type)只能用在switch后面
		case int:
			rect += v // v是一个int类型的变量
		case float32:
			rect += int(v) // v是一个float32类型的变量
		default:
			fmt.Printf("不支持的数据类型 %[1]T 值为 %[1]v\n", ele)
		}
	}
	return rect
}

func main25() {
	// 空接口的典型应用场景
	fmt.Println()
	fmt.Println(1)
	fmt.Println("1")

	rect := SumI(1, float32(3.14), false, "abc")
	fmt.Println(rect)

	//空接口没有定义任何方法，因此任意类型都实现了空接口
	var a any
	a = 3
	a = "ewqrf"
	a = Car{}
	_ = a

	//map的key和value都可以是interface{}类型
	mmap := make(map[interface{}]interface{}, 10)
	mmap["a"] = 1
	mmap["b"] = "A"
	mmap["c"] = 0.2
	mmap[9] = 18
	for k, v := range mmap {
		fmt.Printf("key type %T %v, value %T %v\n", k, k, v, v)
	}
	fmt.Println()

	slice := make([]EmptyInterface, 0, 10)
	slice = append(slice, 1)
	slice = append(slice, "A")
	slice = append(slice, 0.2)
	slice = append(slice, byte(100))                   //go语言中byte是uint8的别名
	fmt.Printf("sum of slice is %d\n", SumI(slice...)) //把切片转为不定长参数
	fmt.Println()

	//类型断言
	var i any
	if v, ok := i.(int); ok { //如果断言成功，则ok为true，v是具体的类型
		fmt.Printf("i是int类型，其值为%d\n", v)
	} else {
		fmt.Println("i不是int类型")
	}
	if v, ok := i.(float32); ok {
		fmt.Printf("i是float类型，其值为%f\n", v)
	} else {
		fmt.Println("i不是float类型")
	}
	//当要判断的类型比较多时，就需要写很多if-else，更好的方法是使用switch i.(type)
	fmt.Println()
}
