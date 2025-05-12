package main

func Foo(a func(int, string) (string, error)) (func(int, string) (string, error), func(int, string) (string, error)) {
	return nil, nil
}

// type功能一：类型别名，缩短名称，少敲几次键盘
type FT func(int, string) (string, error)

// 如果只是想搞类型别名，就用=
// type FT=func(int, string) (string, error)

func Goo(a FT) (FT, FT) {
	return nil, nil
}

type Stack []int

// type功能二：添加自定义方法
func (s *Stack) Push(ele int) {
	*s = append(*s, ele)
}

// type功能三：不需要struct，实现面向接口编程
type Collection interface {
	Push(ele int)
}

type rune int

func main1() {
	// s := new(stack)
	// s := make(Stack, 0, 10)
	s := &Stack{}
	s.Push(1)

	var c Collection = s
	c.Push(2)

	a := new(FT)
	_ = a

	var r rune
	var i int = int(r)
	r = rune(i)
	_ = r
}
