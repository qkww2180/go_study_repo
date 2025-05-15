package main

import (
	"bufio"
	"dqq/io/util"
	"os"
	"testing"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	TEST_LOG = "天对地，雨对风。大陆对长空。山花对海树，赤日对苍穹。雷隐隐，雾蒙蒙。日下对天中。风高秋月白，雨霁晚霞红。牛女二星河左右，参商两曜斗西东。十月塞边，飒飒寒霜惊戍旅；三冬江上，漫漫朔雪冷鱼翁。"
)

var (
	LogRus *logrus.Logger
)

func InitLogrus(logFile string, logLevel logrus.Level) {
	logFile = util.RootPath + "log/" + logFile
	LogRus = logrus.New()
	LogRus.SetLevel(logLevel)
	LogRus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", // 显示ms
	})
	fout, err := rotatelogs.New(
		logFile+".%Y%m%d%H",                      //指定日志文件的路径和名称
		rotatelogs.WithLinkName(logFile),         //为最新的一份日志创建软链接，在windows上就是一个快捷方式（并不是把日志写入了2个文件）
		rotatelogs.WithRotationTime(1*time.Hour), //每隔1小时生成一份新的日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),    //只留最近7天的日志，或使用WithRotationCount只保留最近的几份日志
	)
	if err != nil {
		panic(err)
	}
	LogRus.SetOutput(fout)       //设置日志文件
	LogRus.SetReportCaller(true) //输出是从哪里调起的日志打印，这样日志里会多出func和file
}

// 测试logrus
func TestLogrus(t *testing.T) {
	InitLogrus("rus.log", logrus.DebugLevel)

	file := "go.mod"
	fin, err := os.Open(util.RootPath + "go.mod")
	if err != nil {
		LogRus.Errorf("open file %s failed: %s", file, err)
		return
	}
	reader := bufio.NewReader(fin)
	lineCnt := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			LogRus.Errorf("read file %s failed: %s", file, err)
			break
		}
		LogRus.Debugf("line length %d", len(line))
		lineCnt++
	}
	LogRus.Infof("total read %d line from file %s", lineCnt, file)
}

func BenchmarkLogrus(b *testing.B) {
	InitLogrus("rus_bench.log", logrus.DebugLevel)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogRus.Info(TEST_LOG)
	}
}

// go test -v ./j_io/logger/logrus_test.go -run=^TestLogrus$ -count=1
// go test ./j_io/logger -bench=^BenchmarkLogrus$ -run=^$ -count=1 -benchmem -benchtime=2s
/**
BenchmarkLogrus-8          25557            107739 ns/op            2715 B/op         30 allocs/op
*/
