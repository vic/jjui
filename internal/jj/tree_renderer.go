package jj

import (
	"strings"
)

type TreeRenderer struct {
	dag  *Dag
	rows []RenderContext
}

func NewTreeRenderer(dag *Dag) *TreeRenderer {
	return &TreeRenderer{
		dag:  dag,
		rows: make([]RenderContext, 0),
	}
}

func (t *TreeRenderer) Update(nodeRenderer TreeNodeRenderer) {
	root := t.dag.GetRoot()
	if root == nil {
		return
	}
	edge := IndirectEdge
	if root.Commit.IsRoot() {
		edge = DirectEdge
	}
	t.rows = make([]RenderContext, 0)
	t.renderNode(0, 0, root, edge, nodeRenderer)
}

func (t *TreeRenderer) View(selectedRevision string, height int, nodeRenderer TreeNodeRenderer) string {
	var buffer strings.Builder
	for _, context := range t.rows {
		for _, line := range parseLines(context.Before) {
			buffer.WriteString(strings.Repeat("│ ", context.Level))
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}

		lines := parseLines(context.Content)

		for i, line := range lines {
			buffer.WriteString(strings.Repeat("│ ", context.ParentLevel))
			if i == 0 {
				buffer.WriteString(strings.Repeat("│ ", context.Level-context.ParentLevel))
				buffer.WriteString(context.Glyph)
				buffer.WriteString("  ")
			} else if i < len(lines)-1 {
				buffer.WriteString(strings.Repeat("│ ", context.Level-context.ParentLevel))
				buffer.WriteString("│  ")
			} else {
				if context.Level > context.ParentLevel {
					buffer.WriteString("├─╯  ")
				} else {
					buffer.WriteString("│  ")
				}
			}
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
		if len(lines) == 1 && context.Level > context.ParentLevel {
			buffer.WriteString(strings.Repeat("│ ", context.ParentLevel))
			buffer.WriteString("├─╯\n")
		}
		if context.EdgeType == IndirectEdge {
			buffer.WriteString(nodeRenderer.RenderElidedRevisions())
			buffer.WriteString("\n")
		}
	}
	return buffer.String()

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
		EdgeType:    edgeType,
	}
	renderer.RenderCommit(node.Commit, &context)
	context.height = len(parseLines(context.Content)) + len(parseLines(context.Before))
	t.rows = append(t.rows, context)
}
