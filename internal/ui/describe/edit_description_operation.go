package describe

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
)

type EditDescriptionOperation struct {
	Overlay  tea.Model
	selected *jj.Commit
}

func (e EditDescriptionOperation) Init() tea.Cmd {
	return e.Overlay.Init()
}

func (e EditDescriptionOperation) Update(msg tea.Msg) (common.Operation, tea.Cmd) {
	var cmd tea.Cmd
	e.Overlay, cmd = e.Overlay.Update(msg)
	return EditDescriptionOperation{Overlay: e.Overlay}, cmd
}

func (e EditDescriptionOperation) Render() string {
	return e.Overlay.View()
}

func (e EditDescriptionOperation) RenderPosition() common.RenderPosition {
	return common.RenderPositionDescription
}

func NewEditDescriptionOperation(commands common.UICommands, selected *jj.Commit, width int) (common.Operation, tea.Cmd) {
	op := EditDescriptionOperation{
		selected: selected,
		Overlay:  New(commands, selected.GetChangeId(), selected.Description, width),
	}
	return op, op.Overlay.Init()
}
