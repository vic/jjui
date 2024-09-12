package dag

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"jjui/internal/jj"
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

var DefaultPalette = Palette{
	CommitShortStyle: commitShortStyle,
	CommitIdRestStyle: commitIdRestStyle,
	AuthorStyle: authorStyle,
	Normal: normal,
}

var HighlightedPalette = Palette{
	CommitShortStyle: lipgloss.NewStyle().Background(highlightColor).Inherit(commitShortStyle),
	CommitIdRestStyle: lipgloss.NewStyle().Background(highlightColor).Inherit(commitIdRestStyle),
	AuthorStyle: lipgloss.NewStyle().Background(highlightColor).Inherit(authorStyle),
	Normal: lipgloss.NewStyle().Background(highlightColor).Inherit(normal),
}

type Palette struct {
	CommitShortStyle lipgloss.Style
	CommitIdRestStyle lipgloss.Style
	AuthorStyle lipgloss.Style
	Normal lipgloss.Style
}

func DefaultRenderer(w io.Writer, row *GraphRow, palette Palette) {
	indent := strings.Repeat("│ ", row.Level)
	glyph := "│ "
	nodeGlyph := "○ "
	if !row.IsFirstChild {
		indent = strings.Repeat("│ ", row.Level-1)
		glyph = "├─╯ "
		nodeGlyph = "│ ○ "
	}
	fmt.Print(w, indent)
	fmt.Print(w, nodeGlyph)
	fmt.Print(w, palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
	fmt.Print(w, palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
	fmt.Print(w, " ")
	fmt.Print(w, palette.AuthorStyle.Render(row.Commit.Author))
	fmt.Println(w)
	// description line
	fmt.Print(w, indent)
	fmt.Print(w, glyph)
	if row.Commit.Description == "" {
		fmt.Println(w, palette.Normal.Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("(no description)"))
	} else {
		fmt.Println(w, palette.Normal.Render(row.Commit.Description))
	}
	if row.Elided {
		fmt.Print(w, indent)
		fmt.Println(w, palette.CommitIdRestStyle.Render("~ (elided revisions)"))
	}
}

type GraphRow struct {
	Commit       *jj.Commit
	Level        int
	IsFirstChild bool
	Elided       bool
}

func BuildGraphRows(root *Node) []GraphRow {
	rows := make([]GraphRow, 0)
	Walk(root, func(node *Node, context RenderContext) {
		rows = append(rows, GraphRow{Commit: node.Commit, Level: context.Level, IsFirstChild: context.IsFirstChild, Elided: context.Elided})
	}, RenderContext{Level: 0, IsFirstChild: true})
	return rows
}
