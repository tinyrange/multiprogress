package multiprogress

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type progress struct {
	startTime   time.Time
	endTime     time.Time
	lastUpdated time.Time
	current     int64
	max         int64
	description string
}

// LastUpdated implements Progress.
func (p *progress) LastUpdated() time.Time {
	return p.lastUpdated
}

// Children implements Progress.
func (p *progress) Children() []RenderTree {
	return []RenderTree{}
}

// Height implements Progress.
func (p *progress) Height() int {
	return 1
}

// SetDescription implements Progress.
func (p *progress) SetDescription(s string) {
	p.description = s
}

// Add implements Progress.
func (p *progress) Add(count int) {
	p.current += int64(count)
	p.lastUpdated = time.Now()
}

// Close implements Progress.
func (p *progress) Close() error {
	p.current = p.max
	p.endTime = time.Now()

	return nil
}

// Write implements Progress.
func (pro *progress) Write(p []byte) (n int, err error) {
	pro.Add(len(p))
	return len(p), nil
}

func formatTime(secs float64) string {
	return fmt.Sprintf("%.04fs", secs)
}

func (p *progress) countProgress() string {
	maxCountLen := len(strconv.Itoa(int(p.max)))

	currentCount := fmt.Sprintf("%s%d", strings.Repeat(" ", maxCountLen-len(strconv.Itoa(int(p.current)))), p.current)

	return fmt.Sprintf("%s/%d", currentCount, p.max)
}

func (p *progress) Render(width int) string {
	// If this progress bar is already completed then return a completed string.
	if !p.endTime.IsZero() {
		completedString := fmt.Sprintf(" completed in %s", p.endTime.Sub(p.startTime).String())

		width -= len(completedString)

		description := p.description
		if width < len(p.description)+4 {
			description = description[:width-4] + "... "
		}

		return description + completedString
	}

	width -= 7 // width of progress percentage.
	width -= 2 // width between percentage and count.
	width -= 7 // width of end.

	currentProgress := float64(p.current) / float64(p.max)

	currentTime := time.Since(p.startTime)

	remainingTime := "-"

	if currentProgress > 0 && currentProgress <= 1 {
		totalDuration := currentTime.Seconds() / currentProgress

		remainingTime = formatTime(totalDuration - currentTime.Seconds())
	}

	width -= len(remainingTime) // width of remaining time.
	width -= 1                  // end

	countProgress := p.countProgress()

	width -= len(countProgress) // width of count indicator.

	var progress = " "
	var description = p.description
	if width < len(p.description)+4 {
		description = description[:width-4] + "... "
	} else {
		width -= len(p.description)
		width -= 4 // width of space before and after progress bar.

		progressBar := make([]rune, width)

		for i := 0; i < width; i++ {
			if int(float64(width)*currentProgress) > i {
				progressBar[i] = 'x'
			} else {
				progressBar[i] = ' '
			}
		}

		progress = " [" + string(progressBar) + "] "
	}

	return fmt.Sprintf("%s%s%06.2f%% [%s : est %s]",
		description,
		progress,
		currentProgress*100,
		countProgress,
		remainingTime,
	)
}

var (
	_ Progress = &progress{}
)

func NewProgress(max int64, description string) Progress {
	return &progress{
		startTime:   time.Now(),
		max:         max,
		description: description,
	}
}
