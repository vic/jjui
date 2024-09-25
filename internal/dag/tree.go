package dag

import (
	"sort"

	"jjui/internal/jj"
)

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
	Commit  *jj.Commit
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

func Build(commits []jj.Commit, parents map[string]string) *Node {
	tree := NewDag()
	for _, commit := range commits {
		tree.AddNode(&commit)
	}

	for _, commit := range commits {
		node := tree.GetNode(&commit)
		for _, parent := range commit.Parents {
			if parentNode := tree.GetNodeByChangeId(parent); parentNode != nil {
				parentNode.AddEdge(node, DirectEdge)
			} else {
				current := parent
				for {
					if p, ok := parents[current]; ok {
						if pn := tree.GetNodeByChangeId(p); pn != nil {
							pn.AddEdge(node, IndirectEdge)
							parents[parent] = current
							break
						}
						current = p
						continue
					}
					break
				}
			}
		}
	}

	root := tree.GetRoot()
	root.CalculateDepth()
	return root
}

func (n *Node) CalculateDepth() {
	var maxDepth int
	for _, children := range n.Edges {
		children.To.CalculateDepth()
		if children.To.Depth > maxDepth {
			maxDepth = children.To.Depth
		}
	}
	n.Depth = maxDepth + 1
}

func (d *Dag) AddNode(c *jj.Commit) (node *Node) {
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

func (d *Dag) GetNode(c *jj.Commit) *Node {
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

func Walk(node *Node, renderer Renderer, context RenderContext) {
	sort.Slice(node.Edges, func(a, b int) bool {
		f := node.Edges[a]
		s := node.Edges[b]
		return f.To.Depth > s.To.Depth
	})
	for i, edge := range node.Edges {
		index := i
		nl := context.Level + 1
		if i == 0 {
			nl = context.Level
		}
		Walk(edge.To, renderer, RenderContext{
			Level:        nl,
			Elided:       edge.Type == IndirectEdge,
			IsFirstChild: index == 0 && len(node.Edges) > 1,
		})
	}
	renderer(node, context)
}
