package jj

import (
	"bufio"
	"strings"
)

type RenderContext struct {
	ParentLevel int
	Level       int
	EdgeType    int
	Glyph       string
	Content     string
	Before      string
	After       string
	height      int
}

type TreeNodeRenderer interface {
	RenderCommit(commit *Commit, context *RenderContext)
	RenderElidedRevisions() string
}

func parseLines(content string) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
