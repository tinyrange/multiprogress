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

type RenderGroup interface {
	RenderTree

	Add(child RenderTree)
}

type Progress interface {
	RenderTree

	io.Writer
	io.Closer

	Add(count int)
	SetDescription(s string)
}

type Logger interface {
	RenderTree

	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
