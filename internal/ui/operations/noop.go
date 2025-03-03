package operations

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/idursun/jjui/internal/ui/common"
)

type Noop struct {
	keyMap common.KeyMappings[key.Binding]
}

func (n *Noop) ShortHelp() []key.Binding {
	return []key.Binding{n.keyMap.Up, n.keyMap.Down, n.keyMap.Quit, n.keyMap.Help, n.keyMap.Refresh, n.keyMap.Preview.Mode, n.keyMap.Revset, n.keyMap.Details.Mode, n.keyMap.Evolog, n.keyMap.Rebase.Mode, n.keyMap.Squash.Mode, n.keyMap.Bookmark.Mode, n.keyMap.Git.Mode}
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

func (n *Noop) Name() string {
	return "normal"
}

func Default(c common.AppContext) *Noop {
	return &Noop{
		keyMap: c.KeyMap(),
	}
}
