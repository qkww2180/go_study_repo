package main

import (
	"fmt"
	"math"
	"strings"
)

var scale float64 = 'z' - 'a' + 1 //26进制

func colCount(maxColName string) (int, error) {
	maxColName = strings.ToLower(maxColName)
	var number int
	for i := len(maxColName) - 1; i >= 0; i-- {
		if maxColName[i] >= 'a' && maxColName[i] <= 'z' {
			number += int(math.Pow(scale, float64(len(maxColName)-1-i))) * int(maxColName[i]-'a'+1)
		} else {
			return 0, fmt.Errorf("invalid column name %s", maxColName)
		}
	}
	return number, nil
}

func main12() {
	fmt.Println(colCount("B"))
	fmt.Println(colCount("aA"))
	fmt.Println(colCount("aZ"))
	fmt.Println(colCount("ZzZ"))
	fmt.Println(colCount("AFZ")) //858
	fmt.Println(colCount("XFD")) //16384
}
