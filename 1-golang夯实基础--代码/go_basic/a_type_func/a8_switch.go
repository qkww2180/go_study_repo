package main

import (
	"fmt"
	"strings"
)

func switch_basic() {
	color := "yellow"
	//用switch-case-default模拟if-else
	switch color {
	case "green": //相当于if color=="green"
		fmt.Println("go")
	default:
		fmt.Println("stop")
	}

	//模拟if-else if
	color = "green"
	switch color { //switch后跟一个变量
	case "green": //case后跟一个常量，只要switch和case后面的值类型相同就行
		fmt.Println("go")
	case "red": //相当于else if
		fmt.Println("stop")
	case "yellow":
		fmt.Println("stop")
	default:
		fmt.Println("以上都不是")
	}

	//模拟if-else if-else
	color = "black"
	switch color {
	case "green":
		fmt.Println("go")
	case "red":
		fmt.Println("stop")
	case "yellow":
		fmt.Println("stop")
	default:
		fmt.Printf("invalid traffic signal: %s\n", strings.ToUpper(color))
	}
}

func add(a int) int {
	return a + 10
}

func switch_expression() {
	var a int = 5
	switch add(a) { //switch后跟一个函数表达式
	case 15: //case后跟一个常量
		fmt.Println("right")
	default:
		fmt.Println("wrong")
	}

	const B = 15
	switch B { //switch后跟一个常量
	case add(a): //case后跟一个函数表达式
		fmt.Println("right")
	default:
		fmt.Println("wrong")
	}
}

func switch_condition() {
	color := "yellow"
	switch color {
	case "green":
		fmt.Println("go")
	case "red", "yellow": //用逗号分隔多个condition，它们之间是“或”的关系，只需要有一个condition满足就行
		fmt.Println("stop")
	}

	//switch后带表达式时，switch-case只能模拟相等的情况；如果switch后不带表达式，case后就可以跟任意的条件表达式
	switch {
	case add(5) > 10:
		fmt.Println("right")
	default:
		fmt.Println("wrong")
	}
}

// fallthrough 强制执行下一个case（或default）
func fall_through(age int) {
	fmt.Printf("您的年龄是%d, 您可以：\n", age)
	switch {
	case age > 50:
		fmt.Println("出任国家首脑")
		fallthrough
	case age > 25:
		fmt.Println("生育子女")
		fallthrough
	case age > 22:
		fmt.Println("结婚")
		fallthrough
	case age > 18:
		fmt.Println("开车")
		fallthrough
	case age > 16:
		fmt.Println("参加工作")
	case age > 15:
		fmt.Println("上高中")
		fallthrough
	case age > 3:
		fmt.Println("上幼儿园")
	}
}

func switch_type() {
	var num interface{} = 6.5
	switch num.(type) { //获取interface的具体类型。.(type)只能用在switch后面
	case int:
		fmt.Println("int")
	case float32:
		fmt.Println("float32")
	case float64:
		fmt.Println("float64")
	case byte:
		fmt.Println("byte")
	default:
		fmt.Println("neither")
	}

	switch value := num.(type) { //相当于在每个case内部申明了一个变量value
	case int: //value已被转换为int类型
		fmt.Printf("number is int %d\n", value)
	case float64: //value已被转换为float64类型
		fmt.Printf("number is float64 %f\n", value)
	case byte, string: //如果case后有多个类型，则value还是interface{}类型
		fmt.Printf("number is inerface %v\n", value)
	default:
		fmt.Println("neither")
	}

	//等价形式
	switch num.(type) {
	case int:
		value := num.(int)
		fmt.Printf("number is int %d\n", value)
	case float64:
		value := num.(float64)
		fmt.Printf("number is float64 %f\n", value)
	case byte:
		value := num.(byte)
		fmt.Printf("number is byte %d\n", value)
	default:
		fmt.Println("neither")
	}
}

func main8() {
	switch_basic()
	switch_expression()
	switch_condition()
	fall_through(20)
	switch_type()
}
