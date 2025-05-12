package logger

import (
	"context"
	"testing"
)

func handler(ctx context.Context) {
	Info("hello golang")
	b(ctx)
}

func b(ctx context.Context) {
	Error("hello world")
}

func TestLog(t *testing.T) {
	SetLogFile("test")
	SetLogLevel(InfoLevel)

	ctx := context.WithValue(context.Background(), "traceId", "123456789")
	handler(ctx)

	ctx = context.WithValue(context.Background(), "traceId", "abcdefg")
	handler(ctx)
}

// go test -v ./util/logger -run=^TestLog$ -count=1
