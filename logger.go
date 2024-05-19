package multiprogress

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type logger struct {
	io.Writer
	slog.Handler
}

// Child implements Logger.
func (l *logger) Child(name string) Logger {
	return l
}

// Progress implements Logger.
func (l *logger) Progress(max int64, description string) Progress {
	panic("unimplemented")
}

func (l *logger) log(level slog.Level, msg string, args ...any) {
	var pc uintptr

	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)

	_ = l.Handle(context.Background(), r)
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

func New(target io.Writer) Logger {
	if target == nil {
		target = os.Stderr
	}
	return &logger{
		Writer:  target,
		Handler: slog.NewTextHandler(target, &slog.HandlerOptions{}),
	}
}
