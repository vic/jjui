package describe

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type EditDescriptionOperation struct {
	Overlay  tea.Model
	selected *jj.Commit
}

func (e EditDescriptionOperation) IsEditing() bool {
	return true
}

func (e EditDescriptionOperation) Init() tea.Cmd {
	return e.Overlay.Init()
}

func (e EditDescriptionOperation) Update(msg tea.Msg) (operations.Operation, tea.Cmd) {
	var cmd tea.Cmd
	e.Overlay, cmd = e.Overlay.Update(msg)
	return EditDescriptionOperation{Overlay: e.Overlay}, cmd
}

func (e EditDescriptionOperation) Render() string {
	return e.Overlay.View()
}

func (e EditDescriptionOperation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionDescription
}

func NewOperation(commands common.UICommands, selected *jj.Commit, width int) (operations.Operation, tea.Cmd) {
	op := EditDescriptionOperation{
		selected: selected,
		Overlay:  New(commands, selected.GetChangeId(), selected.Description, width),
	}
	return op, op.Overlay.Init()
}
