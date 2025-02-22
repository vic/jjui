package common

import tea "github.com/charmbracelet/bubbletea"

type EditDescriptionOperation struct {
	Overlay tea.Model
}

func (e EditDescriptionOperation) Render() string {
	return e.Overlay.View()
}

func (e EditDescriptionOperation) RenderPosition() RenderPosition {
	return RenderPositionDescription
}
