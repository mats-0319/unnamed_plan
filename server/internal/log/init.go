package mlog

import (
	"fmt"
	"log"
	"os"
	"time"

	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
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

func Field(msg string, value any) string {
	res := ""

	switch v := value.(type) {
	case error:
		res = fmt.Sprintf("%s - %s", msg, v.Error())
	case string:
		res = fmt.Sprintf("%s: %s", msg, v)
	case bool:
		res = fmt.Sprintf("%s: %t", msg, v)
	case float32, float64:
		res = fmt.Sprintf("%s: %.2f", msg, value)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		res = fmt.Sprintf("%s: %d", msg, v)
	default: // regard as struct
		res = fmt.Sprintf("%s. type: %T, value: %+v", msg, value, value)
	}

	return res
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
	switch mconfig.GetLevel() {
	case "dev":
		level = zap.DebugLevel
	// maybe more levels?
	default:
		level = zap.InfoLevel
	}

	return level
}
