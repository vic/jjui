package dag

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/jj"
	"strings"
)

type RenderContext struct {
	Level        int
	Elided       bool
	IsFirstChild bool
}

type Renderer func(node *Node, context RenderContext)

type GraphRow struct {
	Node         *Node
	Commit       *jj.Commit
	Level        int
	IsFirstChild bool
	Elided       bool
}

func BuildGraphRows(root *Node) []GraphRow {
	rows := make([]GraphRow, 0)
	Walk(root, func(node *Node, context RenderContext) {
		rows = append(rows, GraphRow{Node: node, Commit: node.Commit, Level: context.Level, IsFirstChild: context.IsFirstChild, Elided: context.Elided})
	}, RenderContext{Level: 0, IsFirstChild: true})
	return rows
}

var highlightColor = lipgloss.Color("#44475a")
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

var DefaultPalette = Palette{
	CommitShortStyle:  commitShortStyle,
	CommitIdRestStyle: commitIdRestStyle,
	AuthorStyle:       authorStyle,
	Normal:            normal,
}

var HighlightedPalette = Palette{
	CommitShortStyle:  lipgloss.NewStyle().Background(highlightColor).Inherit(commitShortStyle),
	CommitIdRestStyle: lipgloss.NewStyle().Background(highlightColor).Inherit(commitIdRestStyle),
	AuthorStyle:       lipgloss.NewStyle().Background(highlightColor).Inherit(authorStyle),
	Normal:            lipgloss.NewStyle().Background(highlightColor).Inherit(normal),
}

type Palette struct {
	CommitShortStyle  lipgloss.Style
	CommitIdRestStyle lipgloss.Style
	AuthorStyle       lipgloss.Style
	Normal            lipgloss.Style
}

func DefaultRenderer(w *strings.Builder, row *GraphRow, palette Palette) {
	indent := strings.Repeat("│ ", row.Level)
	glyph := "│"
	nodeGlyph := "○ "
	if !row.IsFirstChild {
		indent = strings.Repeat("│ ", row.Level-1)
		glyph = "├─╯"
		nodeGlyph = "│ ○ "
	}
	w.WriteString(indent)
	w.WriteString(nodeGlyph)
	w.WriteString(palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
	w.WriteString(palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
	w.WriteString(" ")
	w.WriteString(palette.AuthorStyle.Render(row.Commit.Author))
	w.WriteString(" ")
	w.WriteString(fmt.Sprintf("edges: %d ", len(row.Node.Edges)))
	w.WriteString(fmt.Sprintf("level: %d ", row.Level))
	w.WriteString("\n")
	// description line
	w.WriteString(indent)
	w.WriteString(glyph)
	w.WriteString(" ")
	if row.Commit.Description == "" {
		w.WriteString(palette.Normal.Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("(no description)"))
	} else {
		w.WriteString(palette.Normal.Render(row.Commit.Description))
	}
	w.WriteString("\n")
	if row.Elided {
		w.WriteString(indent)
		w.WriteString(palette.CommitIdRestStyle.Render("~ (elided revisions)"))
		w.WriteString("\n")
	}
}
