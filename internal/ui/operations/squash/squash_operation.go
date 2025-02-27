package squash

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	context common.AppContext
	From    string
	Current *jj.Commit
}

var (
	Apply  = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply"))
	Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (s *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Apply):
		return tea.Batch(common.Close, s.context.RunInteractiveCommand(jj.Squash(s.From, s.Current.ChangeIdShort), common.Refresh(s.Current.ChangeIdShort)))
	case key.Matches(msg, Cancel):
		return common.Close
	}
	return nil
}

func (s *Operation) SetSelectedRevision(commit *jj.Commit) {
	s.Current = commit
}

func (s *Operation) Render() string {
	return common.DropStyle.Render("<< into >>")
}

func (s *Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionGlyph
}

func (s *Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		Apply,
		Cancel,
	}
}

func (s *Operation) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func NewOperation(context common.AppContext, from string) *Operation {
	return &Operation{
		context: context,
		From:    from,
	}
}
