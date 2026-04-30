package mlog

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestCustomLog(t *testing.T) {
	Initialize(true)
	defer Close()

	logger := slog.Default().WithGroup("group1").WithGroup("group2")

	Log(logger, slog.LevelDebug, "log msg",
		slog.Int("key1", 10), slog.String("key2", "value2"))
}

func TestLogLevel(t *testing.T) {
	Initialize(true)
	defer Close()

	Debug("debug level log", slog.Any("error", errors.New("debug error")))
	Info("info level log")
	Warn("warn level log")
	Error("error level log", slog.Any("error", errors.New("error")))
}

func TestLogSplitFile(t *testing.T) {
	Initialize(true)
	defer Close()

	lastSize := handler.Size
	currentSize := handler.Size
	for lastSize <= currentSize { // log split
		lastSize = currentSize

		slog.SetDefault(slog.New(handler.WithGroup("groupName")))
		Debug("test log message",
			slog.String("Key1", "Value1"),
			slog.String("Key2", "Value2"))

		currentSize = handler.Size
	}

	time.Sleep(time.Second) // 阻塞，避免异步压缩goroutine因主程序退出而退出
}
