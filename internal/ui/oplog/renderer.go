package oplog

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"strings"
)

type Renderer struct {
	buffer    bytes.Buffer
	lineCount int
	Width     int
}

func (r *Renderer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.lineCount += bytes.Count(p, []byte("\n"))
	return r.buffer.Write(p)
}

func (r *Renderer) LineCount() int {
	return r.lineCount
}

func (r *Renderer) String(start, end int) string {
	lines := strings.Split(r.buffer.String(), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	for end > len(lines) {
		lines = append(lines, "")
	}
	return strings.Join(lines[start:end], "\n")
}

func (r *Renderer) Reset() {
	r.buffer.Reset()
	r.lineCount = 0
}

func (r *Renderer) RenderRow(row Row, highlighted bool) {
	highlightColor := lipgloss.AdaptiveColor{
		Light: config.Current.UI.HighlightLight,
		Dark:  config.Current.UI.HighlightDark,
	}
	highlightSeq := lipgloss.ColorProfile().FromColor(highlightColor).Sequence(true)

	for _, rowLine := range row.Lines {
		lw := strings.Builder{}
		for _, segment := range rowLine.Segments {
			if highlighted {
				fmt.Fprint(&lw, segment.WithBackground(highlightSeq))
			} else {
				fmt.Fprint(&lw, segment.String())
			}
		}
		line := lw.String()
		fmt.Fprint(r, line)
		if highlighted {
			width := lipgloss.Width(line)
			gap := r.Width - width
			if gap > 0 {
				fmt.Fprintf(r, "\033[%sm%s\033[0m", highlightSeq, strings.Repeat(" ", gap))
			}
		}
		fmt.Fprint(r, "\n")
	}
}
