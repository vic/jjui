package dag

import (
	"jjui/internal/jj"
)

type RenderContext struct {
	Level        int
	Elided       bool
	IsFirstChild bool
}

type Renderer func(node *Node, context RenderContext)

type GraphRow struct {
	Node   *Node
	Commit *jj.Commit
	RenderContext
}

func BuildGraphRows(root *Node) []GraphRow {
	rows := make([]GraphRow, 0)
	Walk(root, func(node *Node, context RenderContext) {
		rows = append(rows, GraphRow{Node: node, Commit: node.Commit, RenderContext: context})
	}, RenderContext{Level: 0, IsFirstChild: true})
	return rows
}
