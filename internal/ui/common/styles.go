package common

import "github.com/charmbracelet/lipgloss"

var (
	Black          = lipgloss.Color("0")
	Red            = lipgloss.Color("1")
	Green          = lipgloss.Color("2")
	Yellow         = lipgloss.Color("3")
	Blue           = lipgloss.Color("4")
	Magenta        = lipgloss.Color("5")
	Cyan           = lipgloss.Color("6")
	White          = lipgloss.Color("7")
	IntenseBlack   = lipgloss.Color("8")
	IntenseRed     = lipgloss.Color("9")
	IntenseGreen   = lipgloss.Color("10")
	IntenseYellow  = lipgloss.Color("11")
	IntenseBlue    = lipgloss.Color("12")
	IntenseMagenta = lipgloss.Color("13")
	IntenseCyan    = lipgloss.Color("14")
	IntenseWhite   = lipgloss.Color("15")
)

var ()

var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Blue)

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(IntenseBlack)

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
	ConfirmationText:  lipgloss.NewStyle().Bold(true).Foreground(Magenta),
	Button:            lipgloss.NewStyle().Foreground(White).PaddingLeft(2).PaddingRight(2),
	FocusedButton:     lipgloss.NewStyle().Foreground(IntenseWhite).Background(Blue).PaddingLeft(2).PaddingRight(2),
	ListTitle:         lipgloss.NewStyle().Bold(true).Foreground(White),
	ListItem:          lipgloss.NewStyle().Foreground(Cyan).PaddingLeft(1).PaddingRight(1),
	Added:             lipgloss.NewStyle().Foreground(Green),
	Deleted:           lipgloss.NewStyle().Foreground(Red),
	Modified:          lipgloss.NewStyle().Foreground(Cyan),
	Renamed:           lipgloss.NewStyle().Foreground(Cyan),
	Hint:              lipgloss.NewStyle().Foreground(IntenseBlack).PaddingLeft(1),
	StatusNormal:      lipgloss.NewStyle(),
	StatusSuccess:     lipgloss.NewStyle().Foreground(Green),
	StatusError:       lipgloss.NewStyle().Foreground(Red),
	StatusMode:        lipgloss.NewStyle().Foreground(Black).Background(IntenseBlue),
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
	ConfirmationText  lipgloss.Style
	Button            lipgloss.Style
	FocusedButton     lipgloss.Style
	ListTitle         lipgloss.Style
	ListItem          lipgloss.Style
	Added             lipgloss.Style
	Deleted           lipgloss.Style
	Modified          lipgloss.Style
	Renamed           lipgloss.Style
	Hint              lipgloss.Style
	StatusNormal      lipgloss.Style
	StatusMode        lipgloss.Style
	StatusSuccess     lipgloss.Style
	StatusError       lipgloss.Style
}
