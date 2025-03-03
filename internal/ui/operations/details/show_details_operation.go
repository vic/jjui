package details

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	Overlay tea.Model
	keyMap  common.KeyMappings[key.Binding]
}

func (s Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		s.keyMap.Up,
		s.keyMap.Down,
		s.keyMap.Cancel,
		s.keyMap.Details.Diff,
		s.keyMap.Details.ToggleSelect,
		s.keyMap.Details.Split,
		s.keyMap.Details.Restore,
	}
}

func (s Operation) FullHelp() [][]key.Binding {
	return [][]key.Binding{s.ShortHelp()}
}

func (s Operation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	s.Overlay, cmd = s.Overlay.Update(msg)
	return Operation{Overlay: s.Overlay}, cmd
}

func (s Operation) Render() string {
	return s.Overlay.View()
}

func (s Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionAfter
}

func (s Operation) Name() string {
	return "details"
}

func NewOperation(context common.AppContext, selected *jj.Commit) (operations.Operation, tea.Cmd) {
	op := Operation{
		Overlay: New(context, selected.ChangeId),
		keyMap:  context.KeyMap(),
	}
	return op, op.Overlay.Init()
}
