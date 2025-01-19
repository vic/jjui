package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/ui/common"
	"strings"
)

var (
	Added    = lipgloss.NewStyle().Foreground(common.Green)
	Deleted  = lipgloss.NewStyle().Foreground(common.Red)
	Modified = lipgloss.NewStyle().Foreground(common.Cyan)
)

type Model struct {
	revision string
	files    []string
	common.Commands
}

func New(revision string, commands common.Commands) tea.Model {
	return Model{
		revision: revision,
		Commands: commands,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Status(m.revision)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "h":
			return m, common.Close
		}
	case common.UpdateCommitStatusMsg:
		m.files = msg
	}
	return m, nil
}

func (m Model) View() string {
	var w strings.Builder
	for _, status := range m.files {
		style := Modified
		if strings.HasPrefix(status, "M") {
			style = Modified
		} else if strings.HasPrefix(status, "D") {
			style = Deleted
		} else {
			style = Added
		}
		w.WriteString(style.Render(status))
		w.WriteString("\n")
	}
	return w.String()
}
