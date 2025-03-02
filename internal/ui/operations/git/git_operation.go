package git

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	context common.AppContext
}

func (o *Operation) IsFocused() bool {
	return true
}

func (o *Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionNil
}

func (o *Operation) Render() string {
	return ""
}

func (o *Operation) Name() string {
	return "Git"
}

var (
	Fetch  = key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "git fetch"))
	Push   = key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "git push"))
	Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (o *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Fetch):
		return o.context.RunCommand(jj.GitFetch(), common.Refresh, common.Close)
	case key.Matches(msg, Push):
		return o.context.RunCommand(jj.GitPush(), common.Refresh, common.Close)
	case key.Matches(msg, Cancel):
		return common.Close
	}
	return nil
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

func NewOperation(context common.AppContext) *Operation {
	return &Operation{
		context: context,
	}
}
