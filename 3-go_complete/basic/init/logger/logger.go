package logger

import "fmt"

var Log int

func InitLogger() {
	Log = 9
	fmt.Printf("Log=%d\n", Log)
}
