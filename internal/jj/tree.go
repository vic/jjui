package jj

const (
	DirectEdge   = 1
	IndirectEdge = 2
)

type Dag struct {
	lookup map[string]*Node
	Nodes  []*Node
}

type Node struct {
	Parents []*Node
	Commit  *Commit
	Edges   []*Edge
	Depth   int
}

type Edge struct {
	To   *Node
	Type int
}

func NewDag() *Dag {
	return &Dag{
		lookup: make(map[string]*Node),
		Nodes:  make([]*Node, 0),
	}
}

func (d *Dag) AddNode(c *Commit) (node *Node) {
	node = &Node{
		Commit: c,
		Edges:  make([]*Edge, 0),
	}
	d.Nodes = append(d.Nodes, node)
	d.lookup[c.ChangeId] = node
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

func (d *Dag) GetNode(c *Commit) *Node {
	return d.GetNodeByChangeId(c.ChangeId)
}

func (d *Dag) GetNodeByChangeId(changeId string) *Node {
	return d.lookup[changeId]
}

func (d *Dag) GetRoot() *Node {
	for _, node := range d.Nodes {
		if node.Parents == nil {
			return node
		}
	}
	return nil
}