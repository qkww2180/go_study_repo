package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"slices"

	"golang.org/x/exp/maps"
)

// 内置了min,max,clear函数   len()   append()
func BuiltInFunction() {
	fmt.Printf("最小%d, 最大%d\n", min(1, 4, 7, 5, 2), max(1, 4, 7, 5, 2))
	list := []int{1, 2, 3}
	mp := map[string]string{"姓名": "大乔乔"}
	clear(list) //所以元素置为0
	clear(mp)   //清空
	fmt.Println(list)
	fmt.Println(mp)
}

// structured log
func SLog() {
	slog.Log(context.Background(), slog.LevelInfo, "user info", "姓名", "大乔乔", "年龄", 18)
}

// slices从试验品(golang.org/x/exp/slices)正式转入标准库(slices)
//
// maps还在试验包(golang.org/x/exp/maps)下
func Collection() {
	list := []int{1, 3, 2}
	fmt.Printf("最大值 %d\n", slices.Max(list))
	fmt.Printf("最小值 %d\n", slices.Min(list))
	fmt.Printf("包含%d吗 %t\n", 3, slices.Contains(list, 3))
	list2 := slices.Clone(list) //深拷贝
	fmt.Printf("两个切掉里的元素完全相同吗 %t\n", slices.Equal(list, list2))
	slices.Sort(list) //排序
	index, exists := slices.BinarySearch(list, 3)
	fmt.Printf("存在 %t, 下标 %d\n", exists, index)

	mp := map[string]string{"姓名": "大乔乔", "注册资金": "100万"}
	fmt.Printf("keys %s\n", strings.Join(maps.Keys(mp), "|"))
}

func main() {
	BuiltInFunction()
	SLog()
	Collection()
}

// go run .\basic\21_version\
