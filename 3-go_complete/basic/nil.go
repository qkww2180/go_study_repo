package main

import "fmt"

type A struct{}

func (A) say() string {
	return "hello"
}

type B struct{}

func (B) say() string {
	return "hello"
}

type C struct {
	a A
	b *B
}

func main8() {
	var c C
	fmt.Println(c.a.say())
	fmt.Println(c)
	fmt.Println(c.b)
	fmt.Println(c.b.say())
}
