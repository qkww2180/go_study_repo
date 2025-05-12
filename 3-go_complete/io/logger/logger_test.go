package main

import (
	"bufio"
	"dqq/io/util"
	"dqq/io/util/logger"
	"os"
	"testing"
)

func InitLogger(logFile string, logLevel int) {
	logger.SetLogFile(logFile)   //指定日志文件 daqiaoqiao_golang/log/io.log，log文件夹需要先创建好
	logger.SetLogLevel(logLevel) //指定最低日志级别为debug
}

// 测试自己实现的logger
func TestLogger(t *testing.T) {
	InitLogger("io.log", logger.DebugLevel)

	file := "go.mod"
	fin, err := os.Open(util.RootPath + file)
	if err != nil {
		logger.Error("open file %s failed: %s", file, err)
		return
	}
	reader := bufio.NewReader(fin)
	lineCnt := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			logger.Error("read file %s failed: %s", file, err)
			break
		}
		logger.Debug("line length %d", len(line))
		lineCnt++
	}
	logger.Info("total read %d line from file %s", lineCnt, file)
}

func BenchmarkLogger(b *testing.B) {
	InitLogger("io_bench.log", logger.DebugLevel)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info(TEST_LOG)
	}
}

// go test -v ./io/logger/logger_test.go -run=^TestLogger$ -count=1
// go test ./io/logger -bench=^BenchmarkLogger$ -run=^$ -count=1 -benchmem -benchtime=2s
/**
BenchmarkLogger-8         149134             13536 ns/op            1024 B/op          7 allocs/op
自己实现的logger速度是logrus的8倍。
*/
