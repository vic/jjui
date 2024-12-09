package jj

import (
	"bufio"
	"strings"
)

type TreeRenderer struct {
	buffer strings.Builder
	dag    *Dag
	rows   []RenderContext
}

type RenderContext struct {
	ParentLevel int
	Level       int
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

func NewTreeRenderer(dag *Dag) *TreeRenderer {
	return &TreeRenderer{
		dag:  dag,
		rows: make([]RenderContext, 0),
	}
}

func (t *TreeRenderer) RenderTree(nodeRenderer TreeNodeRenderer) string {
	t.buffer.Reset()
	t.renderNode(0, 0, t.dag.GetRoot(), DirectEdge, nodeRenderer)
	return t.buffer.String()
}

func (t *TreeRenderer) renderNode(level int, parentLevel int, node *Node, edgeType int, renderer TreeNodeRenderer) {
	if node == nil {
		return
	}
	// last edge is to the top node
	for i := len(node.Edges) - 1; i >= 0; i-- {
		edge := node.Edges[i]
		if i == len(node.Edges)-1 {
			t.renderNode(level, level, edge.To, edge.Type, renderer)
		} else {
			t.renderNode(level+1, level, edge.To, edge.Type, renderer)
		}
	}
	context := RenderContext{
		ParentLevel: parentLevel,
		Level:       level,
	}
	renderer.RenderCommit(node.Commit, &context)
	t.rows = append(t.rows, context)
	for _, line := range parseLines(context.Before) {
		t.buffer.WriteString(strings.Repeat("│ ", context.Level))
		t.buffer.WriteString(line)
		t.buffer.WriteString("\n")
		context.height++
	}

	lines := parseLines(context.Content)
	context.height += len(lines)

	for i, line := range lines {
		t.buffer.WriteString(strings.Repeat("│ ", context.ParentLevel))
		if i == 0 {
			t.buffer.WriteString(strings.Repeat("│ ", context.Level-context.ParentLevel))
			t.buffer.WriteString(context.Glyph)
			t.buffer.WriteString("  ")
		} else if i < len(lines)-1 {
			t.buffer.WriteString(strings.Repeat("│ ", context.Level-context.ParentLevel))
			t.buffer.WriteString("│  ")
		} else {
			if level > parentLevel {
				t.buffer.WriteString("├─╯  ")
			} else {
				t.buffer.WriteString("│  ")
			}
		}
		t.buffer.WriteString(line)
		t.buffer.WriteString("\n")
	}
	if len(lines) == 1 && context.Level > context.ParentLevel {
		t.buffer.WriteString(strings.Repeat("│ ", context.ParentLevel))
		t.buffer.WriteString("├─╯\n")
	}
	if edgeType == IndirectEdge {
		t.buffer.WriteString(renderer.RenderElidedRevisions())
		t.buffer.WriteString("\n")
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
