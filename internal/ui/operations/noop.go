package operations

import (
	"github.com/charmbracelet/bubbles/key"
)

type Noop struct{}

var (
	Up            = key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up"))
	Down          = key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down"))
	Details       = key.NewBinding(key.WithKeys("l", "details"), key.WithHelp("l", "details"))
	Apply         = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply"))
	Cancel        = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	Abandon       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "abandon"))
	Edit          = key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
	Diffedit      = key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "diff edit"))
	Split         = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "split"))
	SquashMode    = key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "squash"))
	RebaseMode    = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase"))
	BookmarkMode  = key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "bookmark"))
	GitMode       = key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "git"))
	Description   = key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "description"))
	Evolog        = key.NewBinding(key.WithKeys("O"), key.WithHelp("O", "evolog"))
	Diff          = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff"))
	New           = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new revision"))
	Revset        = key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "edit revset"))
	Refresh       = key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "refresh"))
	Undo          = key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo"))
	Quit          = key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))
	PreviewToggle = key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "toggle preview"))
	Help          = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "show help"))
)

func (n *Noop) ShortHelp() []key.Binding {
	return []key.Binding{Up, Down, Quit, Refresh, PreviewToggle, Revset, Details, Evolog, RebaseMode, SquashMode, BookmarkMode, GitMode, Help}
}

func (n *Noop) FullHelp() [][]key.Binding {
	return [][]key.Binding{n.ShortHelp()}
}

func (n *Noop) RenderPosition() RenderPosition {
	return RenderPositionNil
}

func (n *Noop) Render() string {
	return ""
}
