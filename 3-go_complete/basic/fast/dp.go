package main

import (
	"fmt"
	"time"
)

func InnerProduct(x []float32, y []float32) float32 {
	var rect float32 = 0
	for i := 0; i < len(x); i++ {
		rect += x[i] * y[i]
	}
	return rect
}

func main() {
	const len2 int = 10000000
	crr := []float32{}
	drr := []float32{}
	for i := 0; i < len2; i++ {
		value := float32(i % 10)
		crr = append(crr, value)
		drr = append(drr, value)
	}

	begin := time.Now().UnixMilli()
	prod := InnerProduct(crr, drr)
	end := time.Now().UnixMilli()
	fmt.Printf("蛮力计算内积 %.0f\t\t用时%d毫秒\n", prod, end-begin)
}

// go run basic/fast/dp.go
