package common

import "github.com/charmbracelet/lipgloss"

var (
	Black      = lipgloss.Color("0")
	Cyan       = lipgloss.Color("6")
	Pink       = lipgloss.Color("4")
	Yellow     = lipgloss.Color("3")
	Red        = lipgloss.Color("1")
	Green      = lipgloss.Color("2")
	Comment    = lipgloss.Color("8")
	Foreground = lipgloss.Color("7")
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
	BookmarksStyle:    branchesStyle,
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
	TimestampStyle lipgloss.Style
	BookmarksStyle lipgloss.Style
	ConflictStyle  lipgloss.Style
	Empty             lipgloss.Style
	NonEmpty          lipgloss.Style
	Normal            lipgloss.Style
	Selected          lipgloss.Style
}
