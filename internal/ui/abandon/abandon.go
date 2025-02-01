package abandon

import (
	"strings"

	"jjui/internal/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	message  string
	revision string
	options  []string
	selected int
	common.Commands
}

var selectedStyle = common.DefaultPalette.Normal.
	Background(common.Blue)

var normalStyle = common.DefaultPalette.Normal

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.selected > 0 {
				m.selected--
			}
		case tea.KeyRight:
			if m.selected < len(m.options) {
				m.selected++
			}
		case tea.KeyEnter:
			switch m.options[m.selected] {
			case "Yes":
				return m, m.Abandon(m.revision)
			case "No":
				return m, common.Close
			}
		case tea.KeyEscape:
			return m, common.Close
		default:
		}
	}
	return m, nil
}

func (m Model) View() string {
	p := common.DefaultPalette

	w := strings.Builder{}
	w.WriteString(p.Normal.Render(m.message))
	for i, option := range m.options {
		w.WriteString(" ")
		if i == m.selected {
			w.WriteString(selectedStyle.Render("["))
			w.WriteString(selectedStyle.Width(8).Render(option))
			w.WriteString(selectedStyle.Render("]"))
		} else {
			w.WriteString(normalStyle.Render("["))
			w.WriteString(normalStyle.Width(8).Render(option))
			w.WriteString(normalStyle.Render("]"))
		}
	}
	return w.String()
}

func New(commands common.Commands, revision string) tea.Model {
	return Model{
		message:  "Are you sure you want to abandon this revision?",
		revision: revision,
		options:  []string{"Yes", "No"},
		selected: 0,
		Commands: commands,
	}
}
