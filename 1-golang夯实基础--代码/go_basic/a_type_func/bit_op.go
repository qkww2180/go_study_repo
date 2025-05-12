package main

import "fmt"

func main9() {
	var a uint64 = 200 // 5 = 4+1, 200=128+64+8 11001000
	fmt.Printf("a=%b\n", a)
	binaryFormat(a)
}

func binaryFormat(a uint64) {
	var c uint64 = 1 << 63
	for i := 0; i < 64; i++ {
		if c&a == c {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		c = c >> 1 // 无符号数右移，高位用0填充
	}
}
