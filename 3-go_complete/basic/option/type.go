package option

import "fmt"

// go源码里定义了一些的别名。比如byte是uint8的别名，那么它们在所有场景下使用方式完全相同
type byte = uint8
type rune = int32
type any = interface{}

func Alias() {
	var a uint8
	var b byte = 100
	a = b
	useByteFunc(a)
}

func useByteFunc(a byte) {}

// 自定义别名
type Car struct {
	Power int
}

func (Car) Run() {}

type Boat = Car

// 定义类型
type Duration int64 //go源码定义Duration

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

func Type() {
	var a int64
	const b int64 = 5
	var c Duration = 5
	a = int64(c)                        //通过小括号进行类型转换
	useDurationFunc(Duration(a) * Hour) //通过小括号进行类型转换
	useDurationFunc(Duration(b) * Hour) //通过小括号进行类型转换
	useDurationFunc(5 * Hour)           //字面量会自动执行类型转换，跟其他参与运算的变量类型保持一致
	// useDurationFunc(a * Hour)           //变量和常量则需要显式地执行类型转换
	// useDurationFunc(b * Hour)           //变量和常量则需要显式地执行类型转换
}

func useDurationFunc(a Duration) {}

// 自定义类型
type Plane Car //Plane和Car从概念上讲是两种不同的类型(虽然通过小括号它们可以互相转换)，只不过它们底层的数据类型一样，即Plane也有一个int成员变量Power
// type Plane struct{ Power int } //跟上一行代码等价

func (Plane) Fly() {}

func usePlane(p Plane) {
	p.Fly()
	fmt.Println(p.Power)
	// p.Run()
	// useCar(p)
}

func useCar(c Car) {}

var (
	// 函数也是一种数据类型
	IncFunc = func(c *Car) { c.Power += 1 }
	DecFunc = func(c *Car) { c.Power -= 1 }
)

// 自定义类型
type CarOption func(*Car) // 函数没有成员变量

var (
	IncPower CarOption = func(c *Car) { c.Power += 1 }            //func向type转换时,可以不做显式的类型转换
	DecPower CarOption = CarOption(func(c *Car) { c.Power -= 1 }) //显式的类型转换
	// 上下两种写法都可以，通过等号赋值的就是等价的
	// IncPower CarOption = IncFunc
	// DecPower CarOption = DecFunc
)

// 通过回调函数修改结构体
func ModifyCar(c *Car, functions ...func(*Car)) {
	for _, function := range functions {
		function(c)
	}
}

// 通过Option模式修改结构体
func UseCarOption(c *Car, opts ...CarOption) {
	for _, opt := range opts { //不定长参数当成切片来使用就可以了
		opt(c)
	}
}

func CompareCallbackAndOption() {
	// Option模式和回调函数模式写代码的方式完全相同(Option模式甚至还需要多定义一个类型)，Option模式的优势何在？
	ModifyCar(new(Car), IncFunc, DecFunc)
	ModifyCar(new(Car), IncPower, DecPower) //type直接转为func,不需要显式的类型转换

	UseCarOption(new(Car), IncFunc, DecFunc)
	UseCarOption(new(Car), IncPower, DecPower)
}
