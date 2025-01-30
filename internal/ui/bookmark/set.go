package bookmark

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/ui/common"
)

type SetBookmarkModel struct {
	revision string
	name     textarea.Model
	common.Commands
}

func (m SetBookmarkModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m SetBookmarkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, common.Close
		case "enter":
			return m, m.SetBookmark(m.revision, m.name.Value())
		}
	}
	var cmd tea.Cmd
	m.name, cmd = m.name.Update(msg)
	return m, cmd
}

func (m SetBookmarkModel) View() string {
	return m.name.View()
}

func NewSetBookmark(commands common.Commands, revision string) tea.Model {
	t := textarea.New()
	t.SetValue("")
	t.Focus()
	t.SetWidth(20)
	t.SetHeight(1)
	t.CharLimit = 120
	t.ShowLineNumbers = false
	return SetBookmarkModel{
		name:     t,
		revision: revision,
		Commands: commands,
	}
}
