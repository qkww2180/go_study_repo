package main

import (
	"dqq/io/util"
	"fmt"
	"log"
	"os"
)

var (
	myLogger *log.Logger
)

func init() {
	if logOut, err := os.OpenFile(util.RootPath+"log/dqq.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664); err != nil { //注意是APPEND模式
		panic(err)
	} else {
		myLogger = log.New(logOut, "[dqq] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile) //Llongfile会输出go文件的绝对路径和行号，但是是执行log.Logger.Printf的行号，不是调用Debug()或Info()的行号
	}
}

func Debug(format string, v ...any) {
	myLogger.Printf("DEBUG "+format, v...) //如果不加...，会把v当成一个slice参数来处理，即只能对应到第一个%占位符
}

func Info(format string, v ...any) {
	myLogger.Printf("INFO "+format, v...)
}

func main() {
	log.Printf("%d + %d = %d", 1, 1, 2) //输出到StdOut
	fmt.Printf("%d + %d = %d\n", 1, 1, 2)
	Debug("%d + %d = %d", 2, 2, 4) //输出到日志文件
	Info("%d + %d = %d", 4, 4, 8)
}

// go run ./j_io/logger
