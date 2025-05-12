package main

import "fmt"

//接口是一组行为规范的集合
type Transporter interface { //定义接口。通常接口名以er结尾
	//接口里面只定义方法，不定义变量
	move(src string, dest string) (int, error) //方法名 (参数列表) 返回值列表
	whistle(int) int                           //参数列表和返回值列表里的变量名可以省略。/ˈwɪs(ə)l/
}

type Steamer interface {
	// move(string, string) (int, error)
	// whistle(int) int
	Transporter //通过匿名成员变量来实现“继承”
	displacement() int
}

type Car struct { //定义结构体时无需要显式声明它要实现什么接口
	name  string
	price int
}

//只要结构体拥有接口里声明的所有方法，就称该结构体“实现了接口”
func (car *Car) move(src string, dest string) (int, error) {
	car.price++
	fmt.Printf("从%s到%s用%s运输需要%d元\n", src, dest, car.name, car.price)
	return car.price, nil
}

func move(src string, dest string, car *Car) (int, error) {
	car.price++
	fmt.Printf("从%s到%s用%s运输需要%d元\n", src, dest, car.name, car.price)
	return car.price, nil
}

func (car Car) whistle(n int) int {
	for i := 0; i < n; i++ {
		fmt.Printf("滴 ")
	}
	fmt.Println()
	return n
}

//结构体可以拥有接口之外的方法
func (car Car) GetName() string {
	return car.name
}

type Shiper struct {
	name   string
	price  int
	tonage int
}

func (ship Shiper) move(src string, dest string) (int, error) {
	fmt.Printf("从%s到%s用%s运输需要%d元\n", src, dest, ship.name, ship.price)
	return ship.price, nil
}

func (ship Shiper) whistle(n int) int {
	for i := 0; i < n; i++ {
		fmt.Printf("呜 ")
	}
	fmt.Println()
	return n
}

//一个struct可以同时实现多个接口
func (ship Shiper) displacement() int {
	return ship.tonage
}

// 接口可以作为函数参数，也可以作为结构体的成员变量（策略模式）
func transport(src, dest string, transporter Transporter) error {
	if _, err := transporter.move(src, dest); err != nil {
		return err
	} else {
		transporter.whistle(3)
		return nil
	}
}

func seaTransport(src, dest string, steamer Steamer) error {
	fmt.Printf("排水量%d\n", steamer.displacement())
	if _, err := steamer.move(src, dest); err != nil {
		return err
	} else {
		steamer.whistle(3)
		return nil
	}
}

func main22() {
	car := &Car{"宝马", 100}
	ship := Shiper{"郑和号", 350, 20000}

	var transporter Transporter
	//接口值有两部分组成, 一个指向该接口的具体类型的指针和另外一个指向该具体类型真实数据的指针
	transporter = car
	transporter.whistle(3)
	transporter = ship
	transporter.whistle(3)

	transport("北京", "天津", car)
	transport("北京", "天津", ship)
	fmt.Println()
	seaTransport("北京", "天津", ship)

	// 通过【断言】判断interface是哪种具体的struct
	if _, ok := transporter.(*Car); ok {
		fmt.Println("transporter是*Car")
	} else if _, ok := transporter.(Shiper); ok {
		fmt.Println("transporter是Shiper")
	}
}
