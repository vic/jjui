package jj

import (
	"bufio"
	"io"
	"strings"
)

type renderedCommit struct {
	height  int
	lines   []string
	content string
}

type TreeRenderer struct {
	buffer          strings.Builder
	dag             *Dag
	renderer        TreeNodeRenderer
	renderedCommits map[string]renderedCommit
}

type TreeNodeRenderer interface {
	RenderCommit(commit *Commit) string
	RenderElidedRevisions() string
	RenderGlyph(commit *Commit) string
}

func NewTreeRenderer(dag *Dag, renderer TreeNodeRenderer) *TreeRenderer {
	return &TreeRenderer{
		dag:             dag,
		renderer:        renderer,
		renderedCommits: make(map[string]renderedCommit),
	}
}

func (t *TreeRenderer) RenderTree() string {
	revisions := t.dag.GetRevisions()
	for _, revision := range revisions {
		content := t.renderer.RenderCommit(revision)
		height := 0
		reader := bufio.NewReader(strings.NewReader(content))
		var lines []string
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				height++
				break
			}
			lines = append(lines, string(line))
			height++
		}
		t.renderedCommits[revision.ChangeIdShort] = renderedCommit{
			height:  height,
			lines:   lines,
			content: content,
		}
	}
	t.renderNode(0, t.dag.GetRoot(), DirectEdge, false)
	return t.buffer.String()
}

func (t *TreeRenderer) renderNode(level int, node *Node, edgeType int, indentedChild bool) {
	if node == nil {
		return
	}
	for i := len(node.Edges) - 1; i >= 0; i-- {
		edge := node.Edges[i]
		if i == len(node.Edges)-1 {
			t.renderNode(level , edge.To, edge.Type, false)
		} else {
			t.renderNode(level + 1, edge.To, edge.Type, true)
		}
	}
	indent := strings.Repeat("│ ", level)
	t.buffer.WriteString(indent)
	t.buffer.WriteString(t.renderer.RenderGlyph(node.Commit))
	rc, _ := t.renderedCommits[node.Commit.ChangeIdShort]
	for i, line := range rc.lines {
		if i != 0 && !indentedChild {
			t.buffer.WriteString(indent + "│  ")
		}
		if i == len(rc.lines)-1 && indentedChild {
			indent = strings.Repeat("│ ", level - 1)
			t.buffer.WriteString(indent + "├─╯  ")
		}
		t.buffer.WriteString(line)
		t.buffer.WriteString("\n")
	}
	if edgeType == IndirectEdge {
		t.buffer.WriteString(t.renderer.RenderElidedRevisions())
		t.buffer.WriteString("\n")
	}
}
