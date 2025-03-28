package graph

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/ui/common"
)

type Renderer struct {
	buffer           bytes.Buffer
	skippedLineCount int
	lineCount        int
	Width            int
}

func (r *Renderer) SkipLines(amount int) {
	r.skippedLineCount = r.skippedLineCount + amount
}

func (r *Renderer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	r.lineCount += bytes.Count(p, []byte("\n"))
	return r.buffer.Write(p)
}

func (r *Renderer) LineCount() int {
	return r.lineCount + r.skippedLineCount
}

func (r *Renderer) String(start, end int) string {
	start = start - r.skippedLineCount
	end = end - r.skippedLineCount
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

func RenderRow(r io.Writer, row Row, renderer RowDecorator, highlighted bool, width int) {
	// will render by extending the previous connections
	before := renderer.RenderBefore(row.Commit)
	if before != "" {
		extended := GraphRowLine{}
		if row.Previous != nil {
			extended = row.Previous.Last(Highlightable).Extend(row.Indent)
		}
		lines := strings.Split(before, "\n")
		for _, line := range lines {
			for _, segment := range extended.Segments {
				fmt.Fprint(r, segment.String())
			}
			fmt.Fprintln(r, line)
		}
	}

	highlightColor := lipgloss.AdaptiveColor{
		Light: config.Current.UI.HighlightLight,
		Dark:  config.Current.UI.HighlightDark,
	}

	highlightSeq := lipgloss.ColorProfile().FromColor(highlightColor).Sequence(true)
	var lastLine *GraphRowLine
	for segmentedLine := range row.RowLinesIter(Including(Highlightable)) {
		lastLine = segmentedLine
		lw := strings.Builder{}
		for i, segment := range segmentedLine.Segments {
			if i == segmentedLine.ChangeIdIdx {
				if decoration := renderer.RenderBeforeChangeId(); decoration != "" {
					fmt.Fprint(&lw, decoration)
				}
			}
			if highlighted && i == segmentedLine.CommitIdIdx {
				if decoration := renderer.RenderBeforeCommitId(); decoration != "" {
					fmt.Fprint(&lw, decoration)
				}
			}
			if highlighted {
				fmt.Fprint(&lw, segment.WithBackground(highlightSeq))
			} else {
				fmt.Fprint(&lw, segment.String())
			}
		}
		if segmentedLine.Flags&Revision == Revision && row.IsAffected {
			style := common.DefaultPalette.Dimmed
			if highlighted {
				style = common.DefaultPalette.Dimmed.Background(highlightColor)
			}
			fmt.Fprint(&lw, style.Render(" (affected by last operation)"))
		}
		line := lw.String()
		fmt.Fprint(r, line)
		if highlighted {
			lineWidth := lipgloss.Width(line)
			gap := width - lineWidth
			if gap > 0 {
				fmt.Fprintf(r, "\033[%sm%s\033[0m", highlightSeq, strings.Repeat(" ", gap))
			}
		}
		fmt.Fprint(r, "\n")
	}

	if row.Commit.IsRoot() {
		return
	}

	afterSection := renderer.RenderAfter(row.Commit)
	if afterSection != "" && lastLine != nil {
		extended := lastLine.Extend(row.Indent)
		lines := strings.Split(afterSection, "\n")
		for _, line := range lines {
			for _, segment := range extended.Segments {
				fmt.Fprint(r, segment.String())
			}
			fmt.Fprintln(r, line)
		}
	}

	for segmentedLine := range row.RowLinesIter(Excluding(Highlightable)) {
		for _, segment := range segmentedLine.Segments {
			fmt.Fprint(r, segment.String())
		}
		fmt.Fprint(r, "\n")
	}
}
