package common

import "github.com/charmbracelet/lipgloss"

var (
	Black       = lipgloss.Color("0")
	Red         = lipgloss.Color("1")
	Green       = lipgloss.Color("2")
	Yellow      = lipgloss.Color("3")
	Blue        = lipgloss.Color("4")
	Magenta     = lipgloss.Color("5")
	Cyan        = lipgloss.Color("6")
	White       = lipgloss.Color("7")
	DarkBlack   = lipgloss.Color("8")
	DarkRed     = lipgloss.Color("9")
	DarkGreen   = lipgloss.Color("10")
	DarkYellow  = lipgloss.Color("11")
	DarkBlue    = lipgloss.Color("12")
	DarkMagenta = lipgloss.Color("13")
	DarkCyan    = lipgloss.Color("14")
	DarkWhite   = lipgloss.Color("15")
)

var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Blue)

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(DarkBlack)

var authorStyle = lipgloss.NewStyle().
	Foreground(Yellow)

var timestampStyle = lipgloss.NewStyle().
	Foreground(Cyan)

var branchesStyle = lipgloss.NewStyle().
	Foreground(Blue)

var conflictStyle = lipgloss.NewStyle().
	Foreground(Red)

var normal = lipgloss.NewStyle().
	Foreground(White)

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
	TimestampStyle    lipgloss.Style
	BookmarksStyle    lipgloss.Style
	ConflictStyle     lipgloss.Style
	Empty             lipgloss.Style
	NonEmpty          lipgloss.Style
	Normal            lipgloss.Style
	Selected          lipgloss.Style
}
