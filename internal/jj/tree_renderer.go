package jj

import "strings"

type renderedCommit struct {
	height  int
	content string
}

type TreeRenderer struct {
	buffer          strings.Builder
	renderedCommits map[string]renderedCommit
}

func NewTreeRenderer() *TreeRenderer {
	return &TreeRenderer{
		renderedCommits: make(map[string]renderedCommit),
	}
}

func (t *TreeRenderer) RenderTree(root *Node) string {
	t.renderNode("", root, DirectEdge)
	return t.buffer.String()
}

func (t *TreeRenderer) renderNode(prefix string, node *Node, edgeType int) {
	for i := len(node.Edges) - 1; i >= 0; i-- {
		edge := node.Edges[i]
		if i == len(node.Edges)-1 {
			t.renderNode(prefix, edge.To, edge.Type)
		} else {
			t.renderNode(prefix+"│ ", edge.To, edge.Type)
			t.buffer.WriteString(prefix + "├─╯\n")
		}
	}
	t.buffer.WriteString(prefix)
	if node.Commit.Immutable {
		t.buffer.WriteString("◆  ")
	} else if node.Commit.IsWorkingCopy {
		t.buffer.WriteString("@  ")
	} else {
		t.buffer.WriteString("○  ")
	}
	rendered := t.renderCommit(node.Commit)
	t.renderedCommits[node.Commit.ChangeId] = rendered
	t.buffer.WriteString(rendered.content)
	t.buffer.WriteString("\n")
	if edgeType == IndirectEdge {
		t.buffer.WriteString("~  (elided revisions)\n")
	}
}

func (t *TreeRenderer) renderCommit(commit *Commit) renderedCommit {
	var b strings.Builder
	b.WriteString(commit.ChangeIdShort)
	content := b.String()
	return renderedCommit{height: strings.Count(content, "\n"), content: content}
}
