package multiprogress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type StringRenderer string

// LastUpdated implements RenderTree.
func (s StringRenderer) LastUpdated() time.Time {
	return time.Time{}
}

// Children implements RenderTree.
func (s StringRenderer) Children() []RenderTree {
	return []RenderTree{}
}

// Height implements RenderTree.
func (s StringRenderer) Height() int {
	return 1
}

// Render implements RenderTree.
func (s StringRenderer) Render(width int) string {
	return string(s)[:min(len(s), width)]
}

type ArrayRenderer []RenderTree

// LastUpdated implements RenderTree.
func (a ArrayRenderer) LastUpdated() time.Time {
	var ret time.Time

	for _, child := range a {
		childLast := child.LastUpdated()
		if childLast.After(ret) {
			ret = childLast
		}
	}

	return ret
}

// Children implements RenderTree.
func (a ArrayRenderer) Children() []RenderTree {
	return a
}

// Height implements RenderTree.
func (a ArrayRenderer) Height() int {
	ret := 0

	for _, child := range a {
		ret += child.Height()
	}

	return ret
}

// Render implements RenderTree.
func (a ArrayRenderer) Render(width int) string {
	return ""
}

var (
	_ RenderTree = StringRenderer("")
	_ RenderTree = ArrayRenderer{}
)

type TerminalRenderer struct {
	out           *os.File
	fps           int
	tree          RenderTree
	isTerm        bool
	outFd         int
	ticker        *time.Ticker
	currentHeight int
}

// Close implements io.Closer.
func (term *TerminalRenderer) Close() error {
	term.ticker.Stop()

	return nil
}

func (t *TerminalRenderer) Start() error {
	tickRate := time.Second
	t.outFd = int(t.out.Fd())
	if term.IsTerminal(t.outFd) {
		tickRate /= time.Duration(t.fps)
		t.isTerm = true
	}

	t.ticker = time.NewTicker(tickRate)

	if err := t.render(); err != nil {
		return err
	}

	go func() {
		for range t.ticker.C {
			if err := t.render(); err != nil {
				t.Close()
				break
			}
		}
	}()

	return nil
}

func (term *TerminalRenderer) clearLines(count int) error {
	if _, err := fmt.Fprint(term.out, strings.Repeat("\033[A\033[2K\r", count)); err != nil {
		return err
	}

	return nil
}

func (t *TerminalRenderer) render() error {
	if err := t.clearLines(t.currentHeight); err != nil {
		return err
	}

	var renderTree func(tree RenderTree, depth int) error

	width, height, err := term.GetSize(t.outFd)
	if err != nil {
		return err
	}

	renderedHeight := 0

	renderTree = func(tree RenderTree, depth int) error {
		if height <= 0 {
			return nil
		}

		if depth >= 0 {
			str := tree.Render(width - (depth * 2))

			if len(str) != 0 {
				if _, err := fmt.Fprintf(t.out, "%s%s\n", strings.Repeat("| ", depth), str); err != nil {
					return err
				}

				renderedHeight += 1
			}
		}

		height -= 1

		for _, child := range tree.Children() {
			if err := renderTree(child, depth+1); err != nil {
				return err
			}
		}

		return nil
	}

	if err := renderTree(t.tree, -1); err != nil {
		return err
	}

	t.currentHeight = renderedHeight

	return nil
}

var (
	_ io.Closer = &TerminalRenderer{}
)

func NewTerminalRenderer(out *os.File, fps int, tree RenderTree) *TerminalRenderer {
	return &TerminalRenderer{out: out, fps: fps, tree: tree}
}
