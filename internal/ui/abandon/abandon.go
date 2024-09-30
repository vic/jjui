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
}

var selectedStyle = common.DefaultPalette.Normal.
	Background(common.Pink)

var normalStyle = common.DefaultPalette.Normal

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.selected > 0 {
				m.selected--
			}
		case "right":
			if m.selected < len(m.options) {
				m.selected++
			}
		case "enter":
			switch m.options[m.selected] {
			case "Yes":
				return m, tea.Sequence(common.Close, common.Abandon(m.revision))
			case "No":
				return m, common.Close
			}
		case "esc":
			return m, common.Close
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

func New(revision string) tea.Model {
	return Model{
		message:  "Are you sure you want to abandon this revision?",
		revision: revision,
		options:  []string{"No", "Yes"},
		selected: 0,
	}
}
