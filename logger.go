package multiprogress

import (
	"bytes"
	"context"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

type annotatedMessage struct {
	RenderTree
	addTime time.Time
}

func (msg *annotatedMessage) LastUpdated() time.Time {
	return msg.addTime
}

var (
	_ RenderTree = &annotatedMessage{}
)

type logger struct {
	RenderGroup
}

// Child implements Logger.
func (l *logger) Child(name string) Logger {
	return l
}

func (l *logger) log(level slog.Level, msg string, args ...any) {
	var pc uintptr

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)

	buf := bytes.Buffer{}

	handle := slog.NewTextHandler(&buf, &slog.HandlerOptions{})

	_ = handle.Handle(context.Background(), r)

	str := buf.String()

	str = strings.Trim(str, "\r\n")

	l.Add(&annotatedMessage{RenderTree: StringRenderer(str), addTime: time.Now()})
}

// Info implements Logger.
func (l *logger) Info(msg string, args ...any) {
	l.log(slog.LevelInfo, msg, args...)
}

// Debug implements Logger.
func (l *logger) Debug(msg string, args ...any) {
	l.log(slog.LevelDebug, msg, args...)
}

// Error implements Logger.
func (l *logger) Error(msg string, args ...any) {
	l.log(slog.LevelError, msg, args...)
}

// Warn implements Logger.
func (l *logger) Warn(msg string, args ...any) {
	l.log(slog.LevelWarn, msg, args...)
}

var (
	_ Logger = &logger{}
)

func NewLogger(group RenderGroup) Logger {
	return &logger{RenderGroup: group}
}
