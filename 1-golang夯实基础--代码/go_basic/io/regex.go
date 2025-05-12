package io

import (
	"fmt"
	"regexp"
)

var (
	reg = regexp.MustCompile(`use time (\d+)ms`)
)

func UseRegex() {
	log := "recall use time 38ms, sort use time 20ms"    //这个字符串里reg模式命中两次
	indexs1 := reg.FindAllSubmatchIndex([]byte(log), -1) //-1表示返回所有匹配上reg的地方
	fmt.Println(indexs1)
	indexs2 := reg.FindAllSubmatchIndex([]byte(log), 1) //1表示只需要返回1处(最靠前的1处)匹配上reg的地方
	fmt.Println(indexs2)
	subMatch := indexs1[0]
	begin, end := subMatch[0], subMatch[1]
	fmt.Println(log[begin:end]) //整体匹配上reg的部分
	begin, end = subMatch[2], subMatch[3]
	fmt.Println(log[begin:end]) //匹配上reg中()的部分
}
