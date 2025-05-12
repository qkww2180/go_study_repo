package main

import "fmt"

func main20() {
	a := "8"
	b := '8'
	fmt.Printf("%[1]T %[2]T %[1]s %[2]d\n", a, b)
}
