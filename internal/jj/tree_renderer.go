package jj

import (
	"strings"
)

type TreeRenderer struct {
	buffer   strings.Builder
	dag      *Dag
	renderer TreeNodeRenderer
}

type RenderContext struct {
	Level         int
	IndentedChild bool
	buffer        *strings.Builder
	lines         []string
	glyphAtLine   int
	glyph         string
}

func (rc *RenderContext) RenderLine(line string) {
	rc.lines = append(rc.lines, line)
}

func (rc *RenderContext) Flush() {
	for i := 0; i < rc.glyphAtLine; i++ {
		rc.buffer.WriteString(strings.Repeat("│ ", rc.Level))
		rc.buffer.WriteString(rc.lines[i])
		rc.buffer.WriteString("\n")
	}

    finalGutterWritten := false
	for i := rc.glyphAtLine; i < len(rc.lines); i++ {
		if i == rc.glyphAtLine {
			rc.buffer.WriteString(strings.Repeat("│ ", rc.Level))
			rc.buffer.WriteString(rc.glyph)
			rc.buffer.WriteString("  ")
		} else {
			if rc.IndentedChild && i == len(rc.lines)-1 {
				rc.buffer.WriteString(strings.Repeat("│ ", rc.Level-1))
				rc.buffer.WriteString("├─╯  ")
                finalGutterWritten = true
			} else {
				rc.buffer.WriteString(strings.Repeat("│ ", rc.Level))
                rc.buffer.WriteString("│  ")
			}
		}
		rc.buffer.WriteString(rc.lines[i])
		rc.buffer.WriteString("\n")
	}

    if !finalGutterWritten && rc.IndentedChild {
        rc.buffer.WriteString(strings.Repeat("│ ", rc.Level-1))
        rc.buffer.WriteString("├─╯\n")
    }
    rc.lines = nil
    rc.glyphAtLine = 0
    rc.glyph = ""
}

func (rc *RenderContext) SetGlyph(glyph string) {
	rc.glyph = glyph
	rc.glyphAtLine = len(rc.lines)
}

type TreeNodeRenderer interface {
	RenderCommit(commit *Commit, context *RenderContext)
	RenderElidedRevisions() string
}

func NewTreeRenderer(dag *Dag, renderer TreeNodeRenderer) *TreeRenderer {
	return &TreeRenderer{
		dag:      dag,
		renderer: renderer,
	}
}

func (t *TreeRenderer) NewLine() {
	t.buffer.WriteString("\n")
}

func (t *TreeRenderer) RenderTree() string {
	t.renderNode(0, t.dag.GetRoot(), DirectEdge, false)
	return t.buffer.String()
}

func (t *TreeRenderer) renderNode(level int, node *Node, edgeType int, indentedChild bool) {
	if node == nil {
		return
	}
	// last edge is to the top node
	for i := len(node.Edges) - 1; i >= 0; i-- {
		edge := node.Edges[i]
		if i == len(node.Edges)-1 {
			t.renderNode(level, edge.To, edge.Type, false)
		} else {
			t.renderNode(level+1, edge.To, edge.Type, true)
		}
	}
	context := RenderContext{
		Level:         level,
		IndentedChild: indentedChild,
		buffer:        &t.buffer,
	}
	t.renderer.RenderCommit(node.Commit, &context)
    context.Flush()
	if edgeType == IndirectEdge {
        t.buffer.WriteString(t.renderer.RenderElidedRevisions())
		t.buffer.WriteString("\n")
	}
}
