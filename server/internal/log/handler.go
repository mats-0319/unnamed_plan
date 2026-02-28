package mlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"sync"
)

// log structure: "[time] [level] [file:line func] message | attrs(with group)"
// demo: "[2026-02-14 15:41:37.230] [DEBUG] [log_test.go:24 log.TestLogLevel] debug level log | error=debug error"

type handler struct {
	writer io.Writer

	level  slog.Level
	attrs  []slog.Attr
	groups []string

	mu sync.Mutex
}

func newHandler(writer io.Writer, level slog.Level) *handler {
	return &handler{
		writer: writer,
		level:  level,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// time and level
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")
	levelStr := r.Level.String()

	output := fmt.Sprintf("[%s] [%s] [%s] %s%s\n", timeStr, levelStr, codePosition(), r.Message, h.logAttrs(r))

	_, err := h.writer.Write([]byte(output))

	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) < 1 {
		return h
	}

	newInstance := &handler{
		writer: h.writer,
		level:  h.level,
		attrs:  make([]slog.Attr, len(h.attrs)+len(attrs)),
		groups: make([]string, len(h.groups)),
	}

	copy(newInstance.attrs, h.attrs)
	copy(newInstance.attrs[len(h.attrs):], attrs)
	copy(newInstance.groups, h.groups)

	return newInstance
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newInstance := &handler{
		writer: h.writer,
		level:  h.level,
		attrs:  make([]slog.Attr, len(h.attrs)),
		groups: make([]string, len(h.groups)+1),
	}

	copy(newInstance.attrs, h.attrs)
	copy(newInstance.groups, h.groups)
	newInstance.groups[len(h.groups)] = name

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

func (h *handler) logAttrs(r slog.Record) string {
	length := len(h.attrs) + r.NumAttrs()
	if length < 1 {
		return ""
	}

	attrSlice := make([]string, 0, length)

	for _, attr := range h.attrs {
		attrSlice = append(attrSlice, logAttr(attr, ""))
	}

	r.Attrs(func(a slog.Attr) bool {
		attrSlice = append(attrSlice, logAttrWithGroups(a, h.groups))
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
