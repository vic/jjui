package git

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Commands common.UICommands
	Current  *jj.Commit
}

func (o *Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionNil
}

func (o *Operation) Render() string {
	return ""
}

var (
	Fetch  = key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "git fetch"))
	Push   = key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "git push"))
	Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (o *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Fetch):
		return tea.Batch(o.Commands.GitFetch(o.Current.ChangeIdShort), common.Close)
	case key.Matches(msg, Push):
		return tea.Batch(o.Commands.GitPush(o.Current.ChangeIdShort), common.Close)
	case key.Matches(msg, Cancel):
		return common.Close
	}
	return nil
}

func (o *Operation) SetSelectedRevision(commit *jj.Commit) {
	o.Current = commit
}

func (o *Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		Fetch,
		Push,
		Cancel,
	}
}

func (o *Operation) FullHelp() [][]key.Binding {
	return [][]key.Binding{o.ShortHelp()}
}

func NewOperation(commands common.UICommands) *Operation {
	return &Operation{
		Commands: commands,
	}
}
