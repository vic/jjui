package dag

import (
	"strings"

	"jjui/internal/jj"

	"github.com/charmbracelet/lipgloss"
)

type RenderContext struct {
	Level        int
	Elided       bool
	IsFirstChild bool
}

type Renderer func(node *Node, context RenderContext)

type GraphRow struct {
	Node   *Node
	Commit *jj.Commit
	RenderContext
}

func BuildGraphRows(root *Node) []GraphRow {
	rows := make([]GraphRow, 0)
	Walk(root, func(node *Node, context RenderContext) {
		rows = append(rows, GraphRow{Node: node, Commit: node.Commit, RenderContext: context})
	}, RenderContext{Level: 0, IsFirstChild: true})
	return rows
}

var (
	black      = lipgloss.Color("#000000")
	cyan       = lipgloss.Color("#8be9fd")
	pink       = lipgloss.Color("#ff79c6")
	yellow     = lipgloss.Color("#f1fa8c")
	red        = lipgloss.Color("#ff5555")
	green      = lipgloss.Color("#50fa7b")
	comment    = lipgloss.Color("#6272a4")
	foreground = lipgloss.Color("#f8f8f2")
)

var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(pink)

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(comment)

var authorStyle = lipgloss.NewStyle().
	Foreground(yellow)

var branchesStyle = lipgloss.NewStyle().
	Foreground(pink)

var normal = lipgloss.NewStyle().
	Foreground(foreground)

var selected = lipgloss.NewStyle().
	Foreground(red)

var emptyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(green))

var nonEmptyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color(yellow))

var DropStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(black).
	Background(red)

var DefaultPalette = Palette{
	CommitShortStyle:  commitShortStyle,
	CommitIdRestStyle: commitIdRestStyle,
	AuthorStyle:       authorStyle,
	BranchesStyle:     branchesStyle,
	Empty:             emptyStyle,
	NonEmpty:          nonEmptyStyle,
	Normal:            normal,
	Selected:          selected,
}

type Palette struct {
	CommitShortStyle  lipgloss.Style
	CommitIdRestStyle lipgloss.Style
	AuthorStyle       lipgloss.Style
	BranchesStyle     lipgloss.Style
	Empty             lipgloss.Style
	NonEmpty          lipgloss.Style
	Normal            lipgloss.Style
	Selected          lipgloss.Style
}

func DefaultRenderer(w *strings.Builder, row *GraphRow, palette Palette, highlighted bool) {
	indent := strings.Repeat("│ ", row.Level)
	glyph := "│"
	nodeGlyph := "○ "
	if !row.IsFirstChild {
		indent = strings.Repeat("│ ", row.Level-1)
		glyph = "├─╯"
		nodeGlyph = "│ ○ "
	}
	w.WriteString(indent)
	if highlighted {
		w.WriteString(palette.Selected.Render(nodeGlyph))
	} else {
		w.WriteString(nodeGlyph)
	}
	w.WriteString(palette.CommitShortStyle.Render(row.Commit.ChangeIdShort))
	w.WriteString(palette.CommitIdRestStyle.Render(row.Commit.ChangeId[len(row.Commit.ChangeIdShort):]))
	w.WriteString(" ")
	w.WriteString(palette.AuthorStyle.Render(row.Commit.Author))
	w.WriteString(" ")
	w.WriteString(palette.BranchesStyle.Render(row.Commit.Branches))
	w.WriteString("\n")
	// description line
	w.WriteString(indent)
	w.WriteString(glyph)
	w.WriteString(" ")
	if row.Commit.Description == "" {
		if row.Commit.Empty {
			w.WriteString(palette.Empty.Render("(empty) (no description)"))
		} else {
			w.WriteString(palette.NonEmpty.Render("(no description)"))
		}
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
