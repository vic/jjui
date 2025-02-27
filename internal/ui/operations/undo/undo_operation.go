package undo

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Overlay tea.Model
}

func (o Operation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	o.Overlay, cmd = o.Overlay.Update(msg)
	return Operation{Overlay: o.Overlay}, cmd
}

func (o Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionTop
}

func (o Operation) Render() string {
	return o.Overlay.View()
}

func NewOperation(context common.AppContext) (operations.Operation, tea.Cmd) {
	model := confirmation.New("Are you sure you want to undo last change?")
	model.AddOption("Yes", context.RunCommand(jj.Undo(), common.Refresh, common.Close), key.NewBinding(key.WithKeys("y")))
	model.AddOption("No", common.Close, key.NewBinding(key.WithKeys("n", "esc")))
	return Operation{Overlay: &model}, model.Init()
}
