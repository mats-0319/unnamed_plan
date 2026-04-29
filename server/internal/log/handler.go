package mlog

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

type Handler struct {
	Writer io.Writer

	File     *os.File
	FileName string
	Size     int64
	MaxSize  int64 // single log file max size, unit: MB

	Level  slog.Level
	Attrs  []slog.Attr
	Groups []string
}

var _ slog.Handler = (*Handler)(nil)

func newHandler(fileName string, maxSize int64, level slog.Level) (*Handler, error) {
	h := (&Handler{}).newWriter(fileName)
	h.MaxSize = maxSize << 20
	h.Level = level
	h.Attrs = []slog.Attr{}
	h.Groups = []string{}

	return h, nil
}

func (h *Handler) close() {
	if h != nil {
		_ = h.File.Close()
	}
}

func (h *Handler) newWriter(fileName string) *Handler {
	if h == nil {
		return &Handler{}
	}

	fileIns, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("open log file failed, error:", err)
		return h
	}

	fileInfo, err := os.Stat(fileName) // 写在打开文件后面，因为可能还没有创建日志文件
	if err != nil {
		log.Fatalln("get file info failed, error:", err)
		return h
	}

	h.Writer = io.MultiWriter(fileIns, os.Stdout)
	h.File = fileIns
	h.FileName = fileName
	h.Size = fileInfo.Size()

	return h
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.Level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) (err error) {
	mu.Lock()
	defer mu.Unlock()

	// time and level
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")
	levelStr := r.Level.String()

	output := fmt.Sprintf("[%s] [%s] [%s] %s%s\n", timeStr, levelStr, codePosition(), r.Message, h.logAttrs(r))

	if _, err = h.Writer.Write([]byte(output)); err != nil {
		return
	}

	h.Size += int64(len(output))

	// log rotate, split file
	if h.Size >= h.MaxSize {
		h.close()

		historyFileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02T15-04-05.000Z07:00"))
		if err = os.Rename(h.FileName, historyFileName); err != nil {
			return
		}

		h.newWriter(h.FileName)
	}

	return
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	newInstance := &Handler{
		Writer: h.Writer,
		Level:  h.Level,
		Attrs:  make([]slog.Attr, len(h.Attrs)+len(attrs)),
		Groups: make([]string, len(h.Groups)),
	}

	copy(newInstance.Attrs, h.Attrs)
	copy(newInstance.Attrs[len(h.Attrs):], attrs)
	copy(newInstance.Groups, h.Groups)

	return newInstance
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newInstance := &Handler{
		Writer: h.Writer,
		Level:  h.Level,
		Attrs:  make([]slog.Attr, len(h.Attrs)),
		Groups: make([]string, len(h.Groups)+1),
	}

	copy(newInstance.Attrs, h.Attrs)
	copy(newInstance.Groups, h.Groups)
	newInstance.Groups[len(h.Groups)] = name

	return newInstance
}

func codePosition() string {
	pc := make([]uintptr, 1)
	runtime.Callers(6, pc)

	fs := runtime.CallersFrames(pc)
	f, _ := fs.Next()

	fileName := f.File
	lastIndex := strings.LastIndex(fileName, "/")
	if lastIndex >= 0 {
		index := strings.LastIndex(fileName[:lastIndex], "/")
		if index >= 0 {
			fileName = fileName[index+1:]
		}
	}

	return fmt.Sprintf("%s:%d", fileName, f.Line)
}

func (h *Handler) logAttrs(r slog.Record) string {
	length := len(h.Attrs) + r.NumAttrs()
	if length < 1 {
		return ""
	}

	attrSlice := make([]string, 0, length)

	for _, attr := range h.Attrs {
		attrSlice = append(attrSlice, logAttr(attr, ""))
	}

	r.Attrs(func(a slog.Attr) bool {
		attrSlice = append(attrSlice, logAttrWithGroups(a, h.Groups))
		return true
	})

	return " |" + strings.Join(attrSlice, " ")
}

func logAttrWithGroups(a slog.Attr, groups []string) string {
	var keyPrefix strings.Builder
	for _, group := range groups {
		keyPrefix.WriteString(group + ".")
	}

	return logAttr(a, keyPrefix.String())
}

func logAttr(a slog.Attr, keyPrefix string) string {
	return fmt.Sprintf(" %s=%v", keyPrefix+a.Key, a.Value)
}
