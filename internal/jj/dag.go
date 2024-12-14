package jj

const (
	DirectEdge   = 1
	IndirectEdge = 2
)

type Dag struct {
	Nodes []*Node
}

type Node struct {
	Parents []*Node
	Commit  *Commit
	Edges   []*Edge
}

type Edge struct {
	To   *Node
	Type int
}

func NewDag() Dag {
	return Dag{
		Nodes: make([]*Node, 0),
	}
}

func (d *Dag) AddNode(c *Commit) (node *Node) {
	node = &Node{
		Commit: c,
		Edges:  make([]*Edge, 0),
	}
	d.Nodes = append(d.Nodes, node)
	return node
}

func (n *Node) AddEdge(other *Node, typ int) {
	e := &Edge{
		To:   other,
		Type: typ,
	}
	other.Parents = append(other.Parents, n)
	n.Edges = append(n.Edges, e)
}

func (d *Dag) GetRoot() *Node {
	for _, node := range d.Nodes {
		if node.Parents == nil {
			return node
		}
	}
	return nil
}

func (d *Dag) GetTreeRows() []TreeRow {
	rows := make([]TreeRow, 0)
	root := d.GetRoot()
	if root == nil {
		return rows
	}
	edge := IndirectEdge
	if root.Commit.IsRoot() {
		edge = DirectEdge
	}
	addRow(&rows, 0, 0, root, edge)
	return rows
}

func addRow(rows *[]TreeRow, level int, parentLevel int, node *Node, edgeType int) {
	if node == nil {
		return
	}
	// last edge is connecting to the node that's at the same level
	for i := len(node.Edges) - 1; i >= 0; i-- {
		edge := node.Edges[i]
		if i == len(node.Edges)-1 {
			addRow(rows, level, level, edge.To, edge.Type)
		} else {
			addRow(rows, level+1, level, edge.To, edge.Type)
		}
	}
	context := TreeRow{
		ParentLevel: parentLevel,
		Level:       level,
		EdgeType:    edgeType,
		Commit:      *node.Commit,
	}
	*rows = append(*rows, context)
}
