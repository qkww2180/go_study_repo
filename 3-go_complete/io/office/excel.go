package main

import (
	"dqq/util"
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	fin, err := excelize.OpenFile(util.RootPath + "/data/学生信息表.xlsx")
	if err != nil {
		log.Printf("打开Excel文件失败: %v", err)
		return
	}

	//取得某个单元格的数据(统一转成string)
	cell, err := fin.GetCellValue("一年级", "A2") //Sheet名称，单元格坐标
	if err != nil {
		log.Printf("取不到单元格里的值")
	} else {
		fmt.Println("A2里的值", cell)
	}

	//遍历一个Sheet
	sum := 0.
	count := 0
	rows, err := fin.GetRows("二年级") //rows是一个二维切片，每一行的列数可能不同(Excel表中每一行末尾连续的空白列会被忽略掉)
	if err != nil {
		log.Printf("无法遍历Sheet: %v", err)
		return
	}
	for _, row := range rows {
		if len(row) >= 4 {
			if score, err := strconv.ParseFloat(row[3], 64); err == nil {
				sum += score
				count++
			}
		}
	}
	avgScore := sum / float64(count)

	//写特定的单元格
	cell = "D" + strconv.Itoa(len(rows)+1)
	err = fin.SetCellValue("二年级", cell, avgScore)
	if err != nil {
		log.Printf("写单元格失败: %v", err)
	} else {
		fmt.Printf("向%s单元格写入内容：%f\n", cell, avgScore)
	}

	//保存文件
	fin.Save()
}
