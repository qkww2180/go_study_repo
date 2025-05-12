package main

import "net/http"

//FV是一个函数类型的变量
var FV = func(arg int) {}

//FT是一种类型（FT没有成员变量）
type FT func(arg int)

//类型可以有自己的成员方法
func (ft FT) Hello(arg int) {
	ft(arg)
}

//类型可以实现接口
type IFC interface {
	Hello(arg int)
}

// 参数是函数
func funcCallBack(f func(arg int), arg int) {
	f(arg)
}

// 参数是接口
func funcInterface(i IFC, arg int) {
	i.Hello(arg)
}

type S struct {
}

func main4() {
	FV(3) //函数后面加括号表示调用函数
	_ = S{}
	FT(FV)(3)       //由于FT没有成员变量(不是struct)，创建FT的实例时不能用{}，而需要用()传递一个函数进来。FT的实例当成一个函数来使用
	FT(FV).Hello(3) //FT的实例可以调用成员方法

	/**
	下面2个函数调用执行的逻辑完全是一样的，只不过funcInterface()要求传的是一个type，所以通过适配器把函数FV转成了类型FT
	*/
	funcCallBack(FV, 3)
	funcInterface(FT(FV), 3)

	//用http.Handle实现路由
	http.Handle("/", http.HandlerFunc(Boy)) //http.HandlerFunc类似于FT，它实现了http.Handler接口
}

func Boy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("boy"))
}
