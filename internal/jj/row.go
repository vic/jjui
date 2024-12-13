package jj

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TreeRow struct {
	ParentLevel    int
	Level          int
	EdgeType       int
	Commit         Commit
	Glyph          string
	Content        string
	Before         string
	ElidedRevision string
	height         int
}

type RowRenderer interface {
	Render(row *TreeRow)
}

func RenderRow(w io.Writer, row TreeRow, renderer RowRenderer) {
	renderer.Render(&row)
	for _, line := range parseLines(row.Before) {
		fmt.Fprintf(w, strings.Repeat("│ ", row.Level))
		fmt.Fprintln(w, line)
	}

	lines := parseLines(row.Content)

	for i, line := range lines {
		fmt.Fprintf(w, strings.Repeat("│ ", row.ParentLevel))
		if i == 0 {
			fmt.Fprintf(w, strings.Repeat("│ ", row.Level-row.ParentLevel))
			fmt.Fprintf(w, row.Glyph+"  ")
		} else if i < len(lines)-1 {
			fmt.Fprintf(w, strings.Repeat("│ ", row.Level-row.ParentLevel))
			fmt.Fprintf(w, "│  ")
		} else {
			if row.Level > row.ParentLevel {
				fmt.Fprintf(w, "├─╯  ")
			} else {
				fmt.Fprintf(w, "│  ")
			}
		}
		fmt.Fprintf(w, line+"\n")
	}
	if len(lines) == 1 && row.Level > row.ParentLevel {
		fmt.Fprintf(w, strings.Repeat("│ ", row.ParentLevel))
		fmt.Fprintln(w, "├─╯")
	}
	if row.ElidedRevision != "" {
		fmt.Fprintf(w, strings.Repeat("│ ", row.Level))
		fmt.Fprintln(w, row.ElidedRevision)
	}
}

func parseLines(content string) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
