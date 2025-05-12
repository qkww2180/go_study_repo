package main

type Tesla struct {
	Energy float64
}

func (tesla *Tesla) Run() {
	tesla.Energy -= 40.0
}

func Go2Work(car Car) {
	car.Run()
}

type Car interface {
	Run()
}

func main6() {
	Go2Work(&Tesla{})
}
