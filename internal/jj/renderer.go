package jj

type RenderContext struct {
	Level        int
	Elided       bool
}

type Renderer func(node *Node, context RenderContext)

type GraphRow struct {
	Node   *Node
	Commit *Commit
	RenderContext
}
