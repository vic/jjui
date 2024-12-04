package jj

import "container/list"

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

func (d *Dag) GetRevisions() []*Commit {
	revisions := list.New()
	root := d.GetRoot()
	if root == nil {
		return make([]*Commit, 0)
	}
	var stack []*Node

	stack = append(stack, root)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for i := len(node.Edges) - 1; i >= 0; i-- {
			edge := node.Edges[i]
			stack = append(stack, edge.To)
		}
		revisions.PushBack(node.Commit)
	}
	var ret []*Commit
	for e := revisions.Back(); e != nil; e = e.Prev() {
		ret = append(ret, e.Value.(*Commit))
	}
	return ret
}
