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

var DropStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Black).
	Background(Red)

var DefaultPalette = Palette{
	Normal:           lipgloss.NewStyle(),
	ImmutableNode:    lipgloss.NewStyle().Foreground(IntenseCyan).Bold(true),
	WorkingCopyNode:  lipgloss.NewStyle().Foreground(Green).Bold(true),
	ChangeId:         lipgloss.NewStyle().Foreground(Magenta).Bold(true),
	CommitId:         lipgloss.NewStyle().Foreground(Blue).Bold(true),
	Rest:             lipgloss.NewStyle().Foreground(IntenseBlack),
	Author:           lipgloss.NewStyle().Foreground(Yellow),
	Timestamp:        lipgloss.NewStyle().Foreground(Cyan),
	Bookmarks:        lipgloss.NewStyle().Foreground(Magenta),
	Conflict:         lipgloss.NewStyle().Foreground(Red),
	EmptyPlaceholder: lipgloss.NewStyle().Foreground(Green).Bold(true),
	Placeholder:      lipgloss.NewStyle().Foreground(Yellow).Bold(true),
	Selected:         lipgloss.NewStyle().Foreground(Red),
	Elided:           lipgloss.NewStyle().Foreground(IntenseBlack),
	ConfirmationText: lipgloss.NewStyle().Foreground(Magenta).Bold(true),
	Button:           lipgloss.NewStyle().Foreground(White).PaddingLeft(2).PaddingRight(2),
	FocusedButton:    lipgloss.NewStyle().Foreground(IntenseWhite).Background(Blue).PaddingLeft(2).PaddingRight(2),
	ListTitle:        lipgloss.NewStyle().Foreground(White).Bold(true),
	ListItem:         lipgloss.NewStyle().Foreground(Cyan).PaddingLeft(1).PaddingRight(1),
	Added:            lipgloss.NewStyle().Foreground(Green),
	Deleted:          lipgloss.NewStyle().Foreground(Red),
	Modified:         lipgloss.NewStyle().Foreground(Cyan),
	Renamed:          lipgloss.NewStyle().Foreground(Cyan),
	Hint:             lipgloss.NewStyle().Foreground(IntenseBlack).PaddingLeft(1),
	StatusNormal:     lipgloss.NewStyle(),
	StatusSuccess:    lipgloss.NewStyle().Foreground(Green),
	StatusError:      lipgloss.NewStyle().Foreground(Red),
	StatusMode:       lipgloss.NewStyle().Foreground(Black).Bold(true).Background(Magenta),
}

type Palette struct {
	Normal           lipgloss.Style
	ImmutableNode    lipgloss.Style
	WorkingCopyNode  lipgloss.Style
	ChangeId         lipgloss.Style
	CommitId         lipgloss.Style
	Rest             lipgloss.Style
	Author           lipgloss.Style
	Timestamp        lipgloss.Style
	Bookmarks        lipgloss.Style
	Conflict         lipgloss.Style
	EmptyPlaceholder lipgloss.Style
	Placeholder      lipgloss.Style
	Selected         lipgloss.Style
	ConfirmationText lipgloss.Style
	Button           lipgloss.Style
	FocusedButton    lipgloss.Style
	ListTitle        lipgloss.Style
	ListItem         lipgloss.Style
	Added            lipgloss.Style
	Deleted          lipgloss.Style
	Modified         lipgloss.Style
	Renamed          lipgloss.Style
	Hint             lipgloss.Style
	StatusNormal     lipgloss.Style
	StatusMode       lipgloss.Style
	StatusSuccess    lipgloss.Style
	StatusError      lipgloss.Style
	Elided           lipgloss.Style
}
