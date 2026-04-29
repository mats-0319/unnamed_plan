package mlog

import (
	"log"
	"log/slog"
	"strings"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
)

var handler *Handler

func Initialize(isTestMode ...bool) {
	var err error
	handler, err = newHandler("log.log", 1, getLogLevel(isTestMode...))
	if err != nil {
		log.Fatalln("open log file failed, error:", err)
	}

	slog.SetDefault(slog.New(handler))

	Info("> Config init.") // 不是这里才初始化的，但是只有这里（日志）初始化之后才能使用自定义结构打印这句话
	Info("> Log init.")
}

func Close() {
	handler.close()
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Log 适用于需要通过'WithAttrs'/'WithGroup'定制输出内容、生成新的logger实例的场景
// 例如：http请求、db操作......
//
// 其实可以把slog default logger / slog level复制到这里，在使用过程中就不需要使用slog
// 用法参考测试代码
// 已调整为与默认打印函数使用相同的调用层级数（日志中的代码位置一项，使用'mlog.Info()'/'mlog.Log()'均显示为该函数位置）
func Log(logger *slog.Logger, level slog.Level, msg string, fields ...any) {
	switch level {
	case slog.LevelDebug:
		logger.Debug(msg, fields...)
	case slog.LevelInfo:
		logger.Info(msg, fields...)
	case slog.LevelWarn:
		logger.Warn(msg, fields...)
	case slog.LevelError:
		logger.Error(msg, fields...)
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
