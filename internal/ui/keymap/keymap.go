package keymap

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Keymap struct {
	Current  rune
	Op       operations.Operation
	Bindings map[rune]any
	Up       key.Binding
	Down     key.Binding
	Details  key.Binding
	Cancel   key.Binding
	Apply    key.Binding
}

type BaseLayer struct {
	Abandon      key.Binding
	Edit         key.Binding
	Diffedit     key.Binding
	Split        key.Binding
	GitMode      key.Binding
	RebaseMode   key.Binding
	SquashMode   key.Binding
	BookmarkMode key.Binding
	Description  key.Binding
	Diff         key.Binding
	New          key.Binding
	Revset       key.Binding
	Refresh      key.Binding
	Undo         key.Binding
	Quit         key.Binding
}

type DetailsLayer struct {
	Diff    key.Binding
	Restore key.Binding
	Split   key.Binding
	Mark    key.Binding
}

func NewKeyMap() Keymap {
	bindings := make(map[rune]any)
	bindings[' '] = BaseLayer{
		Abandon:      key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "abandon")),
		Edit:         key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		Diffedit:     key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "diff edit")),
		Split:        key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "split")),
		SquashMode:   key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "squash")),
		RebaseMode:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase")),
		BookmarkMode: key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "bookmark")),
		GitMode:      key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "git")),
		Description:  key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "description")),
		Diff:         key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		New:          key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		Revset:       key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "revset")),
		Refresh:      key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "refresh")),
		Undo:         key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "undo")),
		Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	}

	bindings['d'] = DetailsLayer{
		Diff:    key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "diff")),
		Restore: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restore selected")),
		Split:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "split selected")),
		Mark:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle selection")),
	}

	return Keymap{
		Current:  ' ',
		Bindings: bindings,
		Up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		Down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Details:  key.NewBinding(key.WithKeys("l", "details"), key.WithHelp("l", "details")),
		Apply:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply")),
		Cancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (k *Keymap) ResetMode() {
	k.Current = ' '
}

func (k *Keymap) DetailsMode() {
	k.Current = 'd'
}

func (k *Keymap) ShortHelp() []key.Binding {
	switch b := k.Bindings[k.Current].(type) {
	case BaseLayer:
		return []key.Binding{k.Up, k.Down, b.Revset, b.New, b.Edit, b.Description, b.Diff, b.Abandon, b.Undo, k.Details, b.Split, b.SquashMode, b.Diffedit, b.RebaseMode, b.GitMode, b.BookmarkMode, b.Quit}
	case DetailsLayer:
		return []key.Binding{k.Up, k.Down, b.Mark, b.Diff, b.Restore, b.Split, k.Cancel}
	default:
		if k.Current == 'd' {
			return []key.Binding{k.Up, k.Down, k.Cancel}
		}
		return []key.Binding{}
	}
}

func (k *Keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
