package common

import "github.com/charmbracelet/lipgloss"

var (
	Black      = lipgloss.Color("#000000")
	Cyan       = lipgloss.Color("#8be9fd")
	Pink       = lipgloss.Color("#ff79c6")
	Yellow     = lipgloss.Color("#f1fa8c")
	Red        = lipgloss.Color("#ff5555")
	Green      = lipgloss.Color("#50fa7b")
	Comment    = lipgloss.Color("#6272a4")
	Foreground = lipgloss.Color("#f8f8f2")
)

var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Pink)

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(Comment)

var authorStyle = lipgloss.NewStyle().
	Foreground(Yellow)

var timestampStyle = lipgloss.NewStyle().
	Foreground(Cyan)

var branchesStyle = lipgloss.NewStyle().
	Foreground(Pink)

var conflictStyle = lipgloss.NewStyle().
	Foreground(Red)

var normal = lipgloss.NewStyle().
	Foreground(Foreground)

var selected = lipgloss.NewStyle().
	Foreground(Red)

var emptyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Green)

var nonEmptyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Yellow)

var DropStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Black).
	Background(Red)

var DefaultPalette = Palette{
	CommitShortStyle:  commitShortStyle,
	CommitIdRestStyle: commitIdRestStyle,
	AuthorStyle:       authorStyle,
	TimestampStyle:    timestampStyle,
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
	TimestampStyle    lipgloss.Style
	BranchesStyle     lipgloss.Style
	ConflictStyle     lipgloss.Style
	Empty             lipgloss.Style
	NonEmpty          lipgloss.Style
	Normal            lipgloss.Style
	Selected          lipgloss.Style
}
