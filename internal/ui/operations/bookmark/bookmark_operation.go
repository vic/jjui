package bookmark

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type ChooseBookmarkOperation struct {
	selected *jj.Commit
	Commands common.UICommands
}

func (c *ChooseBookmarkOperation) SetSelectedRevision(commit *jj.Commit) {
	c.selected = commit
}

func (c *ChooseBookmarkOperation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Move):
		//TODO: pass the correct width
		operation, cmd := NewMoveBookmarkOperation(c.Commands, c.selected, 80)
		return tea.Sequence(common.SetOperation(operation), cmd)
	case key.Matches(msg, Set):
		operation, cmd := NewSetBookmarkOperation(c.Commands, c.selected)
		return tea.Sequence(common.SetOperation(operation), cmd)
	case key.Matches(msg, Delete):
		//TODO: pass the correct width
		operation, cmd := NewDeleteBookmarkOperation(c.Commands, c.selected, 80)
		return tea.Sequence(common.SetOperation(operation), cmd)
	case key.Matches(msg, Cancel):
		return common.Close
	}
	return nil
}

var (
	Move   = key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "bookmark move"))
	Set    = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "bookmark set"))
	Delete = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "bookmark delete"))
	Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (c *ChooseBookmarkOperation) ShortHelp() []key.Binding {
	return []key.Binding{
		Move,
		Set,
		Delete,
		Cancel,
	}
}

func (c *ChooseBookmarkOperation) FullHelp() [][]key.Binding {
	return [][]key.Binding{c.ShortHelp()}
}

func (c *ChooseBookmarkOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionNil
}

func (c *ChooseBookmarkOperation) Render() string {
	return ""
}

func NewChooseBookmarkOperation(commands common.UICommands) operations.Operation {
	return &ChooseBookmarkOperation{
		Commands: commands,
	}
}
