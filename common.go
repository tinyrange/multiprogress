package multiprogress

import (
	"io"
	"time"
)

type RenderTree interface {
	Render(width int) string
	Height() int
	Children() []RenderTree
	LastUpdated() time.Time
}

type Progress interface {
	RenderTree

	io.Writer
	io.Closer

	Add(count int)
	SetDescription(s string)
}

type Logger interface {
	io.Writer

	Child(name string) Logger

	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)

	Progress(max int64, description string) Progress
}
