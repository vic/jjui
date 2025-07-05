package oplog

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/ui/common"
	"io"
	"strings"
)

type iterator struct {
	Palette             common.Palette
	HighlightBackground lipgloss.AdaptiveColor
	Width               int
	Rows                []row
	isHighlighted       bool
	current             int
	Cursor              int
}

func newIterator(rows []row, cursor int, width int) *iterator {
	return &iterator{
		Palette:             common.DefaultPalette,
		HighlightBackground: lipgloss.AdaptiveColor{Light: config.Current.UI.HighlightLight, Dark: config.Current.UI.HighlightDark},
		Width:               width,
		Rows:                rows,
		isHighlighted:       false,
		current:             -1,
		Cursor:              cursor,
	}
}

func (o *iterator) IsHighlighted() bool {
	return o.current == o.Cursor
}

func (o *iterator) Render(r io.Writer) {
	row := o.Rows[o.current]
	renderer := o
	highlightColor := lipgloss.AdaptiveColor{
		Light: config.Current.UI.HighlightLight,
		Dark:  config.Current.UI.HighlightDark,
	}
	highlightSeq := lipgloss.ColorProfile().FromColor(highlightColor).Sequence(true)

	for _, rowLine := range row.Lines {
		lw := strings.Builder{}
		for _, segment := range rowLine.Segments {
			if o.isHighlighted {
				fmt.Fprint(&lw, segment.WithBackground(highlightSeq).String())
			} else {
				fmt.Fprint(&lw, segment.String())
			}
		}
		line := lw.String()
		fmt.Fprint(r, line)
		if o.isHighlighted {
			lineWidth := lipgloss.Width(line)
			gap := renderer.Width - lineWidth
			if gap > 0 {
				fmt.Fprintf(r, "\033[%sm%s\033[0m", highlightSeq, strings.Repeat(" ", gap))
			}
		}
		fmt.Fprint(r, "\n")
	}
}

func (o *iterator) RowHeight() int {
	return len(o.Rows[o.current].Lines)
}

func (o *iterator) Next() bool {
	o.current++
	if o.current >= len(o.Rows) {
		return false
	}
	o.isHighlighted = o.current == o.Cursor
	return true
}

func (o *iterator) Len() int {
	return len(o.Rows)
}
