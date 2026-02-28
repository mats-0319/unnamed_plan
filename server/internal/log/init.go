package mlog

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
)

func Initialize(isTestMode ...bool) {
	var err error
	fileIns, err = os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("create log file failed, error:", err)
	}

	multiWriter := io.MultiWriter(fileIns, os.Stdout)
	handlerIns := newHandler(multiWriter, getLogLevel(isTestMode...))

	logger := slog.New(handlerIns)
	slog.SetDefault(logger)

	Info("> Config init.")
	Info("> Log init.")
}

var fileIns *os.File

func Close() {
	if fileIns != nil {
		_ = fileIns.Close()
	}
}

func Debug(msg string, args ...*LogField) {
	slog.Debug(msg, fieldsToAnySlice(args)...)
}

func Info(msg string, args ...*LogField) {
	slog.Info(msg, fieldsToAnySlice(args)...)
}

func Warn(msg string, args ...*LogField) {
	slog.Warn(msg, fieldsToAnySlice(args)...)
}

func Error(msg string, args ...*LogField) {
	slog.Error(msg, fieldsToAnySlice(args)...)
}

type LogField struct {
	Key   string
	Value any
}

func Field(msg string, value any) *LogField {
	return &LogField{
		Key:   msg,
		Value: value,
	}
}

func fieldsToAnySlice(fields []*LogField) []any {
	fs := make([]any, 0, len(fields)*2)
	for _, field := range fields {
		fs = append(fs, field.Key, field.Value)
	}

	return fs
}

// Log 适用于需要通过'WithAttrs'/'WithGroup'定制输出内容、生成新的logger实例的场景
// 例如：http请求、db操作......
//
// 其实可以把slog default logger / slog level复制到这里，在使用过程中就不需要使用slog
// 用法参考测试代码
// 已调整为与默认打印函数使用相同的调用层级数（日志中的代码位置一项，使用'mlog.Info()'/'mlog.Log()'均显示为该函数位置）
func Log(logger *slog.Logger, level slog.Level, msg string, fields ...*LogField) {
	fs := fieldsToAnySlice(fields)

	switch level {
	case slog.LevelDebug:
		logger.Debug(msg, fs...)
	case slog.LevelInfo:
		logger.Info(msg, fs...)
	case slog.LevelWarn:
		logger.Warn(msg, fs...)
	case slog.LevelError:
		logger.Error(msg, fs...)
	default:
		logger.Error("unknown level: " + level.String())
	}
}

func getLogLevel(isTestMode ...bool) slog.Level {
	if len(isTestMode) > 0 && isTestMode[0] {
		return slog.LevelDebug
	}

	levelStr := mconfig.GetLevel()

	var level slog.Level
	switch strings.ToLower(levelStr) {
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	default: // 'debug' and other unknown levels
		level = slog.LevelDebug
	}

	return level
}
