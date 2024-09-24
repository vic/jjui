package common

import "github.com/charmbracelet/lipgloss"

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

var conflictStyle = lipgloss.NewStyle().
	Foreground(red)

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
	ConflictStyle:     conflictStyle,
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
	ConflictStyle     lipgloss.Style
	Empty             lipgloss.Style
	NonEmpty          lipgloss.Style
	Normal            lipgloss.Style
	Selected          lipgloss.Style
}
