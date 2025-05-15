package j_io

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
)

// 定制Logger
func NewLogger(logFile string) *log.Logger {
	//以append方式打开日志文件
	fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open log file failed: %v\n", err)
	}
	logger := log.New(fout, "[MY_BIZ]", log.Ldate|log.Lmicroseconds) //通过flag参数定义日志的格式，时间精确到微秒1E-6s
	return logger
}

func Log(logger *log.Logger) {
	logger.Printf("%d+%d=%d\n", 3, 4, 3+4)
	logger.Println("Hello 大乔乔")
	// logger.Fatalln("Bye, the world") //日志输出后会执行os.Exit(1)

	log.Printf("%d+%d=%d\n", 3, 4, 3+4)
	log.Println("Hello 大乔乔")
	// log.Fatalln("Bye, the world")
}

// 定制SLogger
func NewSLogger(logFile string) *slog.Logger {
	//以append方式打开日志文件
	fout, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open log file failed: %v\n", err)
	}
	logger := slog.New(
		// slog.NewJSONHandler(fout, &slog.HandlerOptions{
		// 	AddSource: true,
		// 	Level:     slog.LevelInfo,
		// }),

		// slog.NewTextHandler(fout, &slog.HandlerOptions{
		// 	AddSource: true,
		// 	Level:     slog.LevelInfo,
		// }),

		&SlogContextHandler{
			slog.NewJSONHandler(fout, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					// check that we are handling the time key
					if a.Key != slog.TimeKey {
						return a
					}

					t := a.Value.Time()

					// change the value from a time.Time to a String
					// where the string has the correct time format.
					a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))

					return a
				},
			}),
		},
	)
	return logger
}

func SLog(logger *slog.Logger) {
	logger.Debug("加法运算", "a", 3, "b", 4, "sum", 3+4)
	logger.Info("加法运算", "a", 3, "b", 4, "sum", 3+4)
	logger.Error("Hello 大乔乔")

	ctx1 := AppendCtx(context.Background(), slog.String("user", "张朝阳"))
	ctx2 := AppendCtx(ctx1, slog.Int("age", 18))
	logger.InfoContext(ctx2, "welcome")

	slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.Debug("加法运算", "a", 3, "b", 4, "sum", 3+4)
	slog.Info("加法运算", "a", 3, "b", 4, "sum", 3+4)
	slog.Error("Hello 大乔乔")
}
