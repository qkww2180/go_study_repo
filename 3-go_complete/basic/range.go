package main

import "fmt"

type User struct {
	Id   int
	Name string
}

func Slice2Map(slice []User) map[int]*User {
	result := make(map[int]*User, len(slice))
	for i, user := range slice {
		// result[i] = &user
		tmp := user
		result[i] = &tmp
		fmt.Printf("%p %+v\n", &tmp, user)
	}
	return result
}

func main3() {
	slice := []User{{Id: 1, Name: "大乔乔"}, {Id: 2, Name: "大脸猫"}}
	mp := Slice2Map(slice)
	for _, v := range mp {
		fmt.Printf("%p %+v\n", v, *v)
	}
}
