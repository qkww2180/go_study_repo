package main

import (
	"fmt"
	"strings"
	"time"
)

func if_expression() {
	if 5 > 9 { //如果逻辑表达式成立，就会执行{}里的内容。逻辑表达式不需要加()
		fmt.Println("5>9")
	}
	if 5 < 9 { //“{”必须紧跟在逻辑表达式后面，不能另起一行
		fmt.Println("5<9")
	}
	var a int = 5
	if a < 9 { //逻辑表达中可以含有变量或常量
		fmt.Println("a<9")
	}
	if b := 8; b < 9 { //if句子中允许包含1个(仅1个)分号，在分号前初始化一些局部变量(即只在if-else块内可见)
		fmt.Println("b<9")
	}
	if c, d, e := 5, 9, 2; c < d && (c > e || c > 3) { //初始化多个局部变量。复杂的逻辑表达式，&&且，||或
		fmt.Println("fit")
	}
}

func else_expression() {
	var color string = "yellow"
	//if-else模式
	if color == "green" {
		fmt.Println("go")
	} else {
		fmt.Println("stop")
	}
	//if-else if模式
	color = "green"
	if color == "red" { //if只能有一个
		fmt.Println("stop")
	} else if color == "green" {
		fmt.Println("go")
	} else if color == "yellow" { //else if可以有一个或者连续多个
		fmt.Println("stop")
	}
	//if-else if-else模式
	color = "black"
	if color == "red" { //if只能有一个
		fmt.Println("stop")
	} else if color == "green" {
		fmt.Println("go")
	} else if color == "yellow" { //else if可以有0个、一个或者连续多个
		fmt.Println("stop")
	} else { //else有0个或1个
		fmt.Printf("invalid traffic signal: %s\n", strings.ToUpper(color))
	}

	m := map[string]bool{"male": true, "female": false}
	if value, exists := m["tree"]; exists { //初始化value和exists这两个局部变量，它们在if和else块里都是可见的
		fmt.Printf("key exists, value is %v\n", value)
	} else {
		fmt.Printf("key NOT exists, value is %v\n", value) //map中如果不存在相应的key，则value是相应类型的默认值
	}
}

// if-else表达式嵌套
func if_nest() {
	now := time.Now()
	fmt.Println(now.Weekday(), now.Hour())
	if now.Weekday().String() != "Saturday" && now.Weekday().String() != "Sunday" {
		if (now.Hour() >= 7 && now.Hour() < 9) || (now.Hour() >= 17 && now.Hour() < 19) {
			fmt.Println("当前为工作日上下班高峰期，公交车道禁止通行")
		} else {
			if now.Hour() >= 21 || now.Hour() <= 5 {
				fmt.Println("夜间请小心驾驶")
			} else {
				fmt.Println("当前公交车道可放心通行")
			}
		}
	} else {
		if now.Hour() >= 21 || now.Hour() <= 5 {
			fmt.Println("夜间请小心驾驶")
		} else {
			fmt.Println("周六日公交车道可正常通行")
		}
	}
}
