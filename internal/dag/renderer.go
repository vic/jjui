package dag

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type RenderContext struct {
	Level        int
	Elided       bool
	IsFirstChild bool
}

type Renderer func(node *Node, context RenderContext)

var highlightColor = lipgloss.Color("#282a36")
var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#bd93f9"))

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#6272a4"))

var authorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ffb86c"))

var normal = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#f8f8f2"))

func DefaultRenderer(node *Node, context RenderContext) {
	indent := strings.Repeat("│ ", context.Level)
	glyph := "│ "
	nodeGlyph := "○ "
	if !context.IsFirstChild {
		indent = strings.Repeat("│ ", context.Level-1)
		glyph = "├─╯ "
		nodeGlyph = "│ ○ "
	}
	fmt.Print(indent)
	fmt.Print(nodeGlyph)
	fmt.Print(commitShortStyle.Render(node.Commit.ChangeIdShort))
	fmt.Print(commitIdRestStyle.Render(node.Commit.ChangeId[len(node.Commit.ChangeIdShort):]))
	fmt.Print(" ")
	fmt.Print(authorStyle.Render(node.Commit.Author))
	fmt.Println()
	// description line
	fmt.Print(indent)
	fmt.Print(glyph)
	if node.Commit.Description == "" {
		fmt.Println(normal.Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("(no description)"))
	} else {
		fmt.Println(normal.Render(node.Commit.Description))
	}
	if context.Elided {
		fmt.Print(indent)
		fmt.Println(commitIdRestStyle.Render("~ (elided revisions)"))
	}
}
