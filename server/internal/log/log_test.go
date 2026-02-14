package mlog

import (
	"errors"
	"log/slog"
	"testing"
)

func TestCustomLog(t *testing.T) {
	Initialize()
	defer Close()

	logger := slog.With("key_1", "value_1").
		WithGroup("groupOne").
		WithGroup("groupTwo")

	Log(logger, slog.LevelDebug, "log msg", Field("key1", 10), Field("error", errors.New("new error")))
}

func TestLogLevel(t *testing.T) {
	Initialize()
	defer Close()

	Debug("debug level log", Field("error", errors.New("debug error")))
	Info("info level log")
	Warn("warn level log")
	Error("error level log", Field("error", errors.New("error")))
}
