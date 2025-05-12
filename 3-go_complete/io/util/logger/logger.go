package logger

import (
	"dqq/io/util"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	//log.Logger写文件支持并发
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	logLevel    = 0 //默认的LogLevel为0，即所有级别的日志都打印

	logFile       string       //日志文件
	logOut        *os.File     //日志文件
	day           int          //当前是一年中的第几天
	dayChangeLock sync.RWMutex //切换日志文件时加锁
)

const (
	DebugLevel = iota //iota=0
	InfoLevel         //InfoLevel=iota, iota=1
	WarnLevel         //WarnLevel=iota, iota=2
	ErrorLevel        //ErrorLevel=iota, iota=3
)

func SetLogLevel(level int) {
	logLevel = level
}

func SetLogFile(file string) {
	logFile = file
	now := time.Now()
	var err error
	//日志文件放到根目录的log子文件夹下，log文件夹需要先创建好
	if logOut, err = os.OpenFile(util.RootPath+"log/"+file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664); err != nil { //注意是APPEND模式
		panic(err)
	} else {
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags) //LstdFlags = Ldate | Ltime
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
		errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
		day = now.YearDay()            //当前是一年中的第几天
		dayChangeLock = sync.RWMutex{} //切换日志文件时加锁
	}
}

// 检查是否需要切换日志文件，如果需要则切换
func checkAndChangeLogfile() {
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()
	now := time.Now()
	if now.YearDay() == day {
		return //不需要切换日志文件
	}

	//关闭老的日志文件
	logOut.Close()

	//给老的日志文件加上日期后缀
	var err error
	postFix := now.Add(-24 * time.Hour).Format("20060102") //昨天的日期
	if err = os.Rename(logFile, logFile+"."+postFix); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("append date postfix %s to log file %s failed: %v\n", postFix, logFile, err)) //如果logger本身出错，则把错误信息打到标准错误输出里
		return
	}

	//打开新的文件日志
	if logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("create log file %s failed %v\n", logFile, err))
		return
	} else {
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		warnLogger = log.New(logOut, "[WARN] ", log.LstdFlags)
		errorLogger = log.New(logOut, "[ERROR] ", log.LstdFlags)
		day = now.YearDay() //更新day
	}
}

// 获取当前运行代码的文件名和行号
func getFileAndLineNo() (string, string, int) {
	funcName, file, line, ok := runtime.Caller(3) //返回函数名、文件名、行号。runtime.Caller(3)是第3层调用堆栈，getFileAndLineNo()第0层 -->  addPrefix()第1层 --> Info()第2层 --> 调用logger.Info()的地方第3层
	if ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	} else {
		return "", "", 0
	}
}

// 给每行日志添加前缀：文件名加行号
func addPrefix() string {
	file, _, line := getFileAndLineNo()
	arr := strings.Split(file, "/")
	if len(arr) > 3 {
		arr = arr[len(arr)-3:] //不需要完整的绝对路径，只取最末3级路径
	}
	return strings.Join(arr, "/") + ":" + strconv.Itoa(line) //文件名加行号
}

func Debug(format string, v ...any) {
	if logLevel <= DebugLevel {
		checkAndChangeLogfile()
		debugLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Info(format string, v ...any) {
	if logLevel <= InfoLevel {
		checkAndChangeLogfile()
		infoLogger.Printf(addPrefix()+" "+format, v...) //format末尾如果没有换行符会自动加上
	}
}

func Warn(format string, v ...any) {
	if logLevel <= WarnLevel {
		checkAndChangeLogfile()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Error(format string, v ...any) {
	if logLevel <= ErrorLevel {
		checkAndChangeLogfile()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	}
}
