package mlog

import (
	"log"
	"os"
	"time"

	mconf "github.com/mats0319/unnamed_plan/server/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zLog *zap.Logger

// Log todo: optimize - log level check
func Log(msg string, fields ...string) {
	fs := make([]zap.Field, len(fields))
	for i := range fields {
		fs[i] = zap.String("", fields[i])
	}

	zLog.Info(msg, fs...)
}

func init() {
	ws, err := logWriteSyncer()
	if err != nil {
		log.Fatalln("init log write syncer failed, error :", err)
	}

	coreSlice := make([]zapcore.Core, 0, 2)
	coreSlice = append(coreSlice, zapcore.NewCore(logEncoder(), ws, logLevel()))
	coreSlice = append(coreSlice, zapcore.NewCore(logEncoder(), os.Stdout, logLevel()))

	core := zapcore.NewTee(coreSlice...)
	zLog = zap.New(core, zap.AddCaller())

	Log("> Config init.")
	Log("> Log init.")
}

func logWriteSyncer() (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("create log file failed, error:", err)
		return nil, err
	}

	return zapcore.AddSync(file), nil
}

func logEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

func logLevel() zapcore.Level {
	var level zapcore.Level
	switch mconf.GetLevel() {
	case "dev":
		level = zap.DebugLevel
	// maybe more levels?
	default:
		level = zap.InfoLevel
	}

	return level
}
