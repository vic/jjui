package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/ui/common"
	"strings"
)

var (
	Added            = lipgloss.NewStyle().Foreground(common.Green)
	AddedSelected    = lipgloss.NewStyle().Inherit(Added).Bold(true).Background(common.DarkBlack)
	Deleted          = lipgloss.NewStyle().Foreground(common.Red)
	DeletedSelected  = lipgloss.NewStyle().Inherit(Deleted).Bold(true).Background(common.DarkBlack)
	Modified         = lipgloss.NewStyle().Foreground(common.Cyan)
	ModifiedSelected = lipgloss.NewStyle().Inherit(Modified).Bold(true).Background(common.DarkBlack)
)

type Model struct {
	revision string
	selected int
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
		case "d":
			return m, m.Commands.GetDiff(m.revision, m.files[m.selected][2:])
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.files)-1 {
				m.selected++
			}
		}
	case common.UpdateCommitStatusMsg:
		m.files = msg
	}
	return m, nil
}

func (m Model) View() string {
	var w strings.Builder
	maxWidth := 0
	for _, status := range m.files {
		if len(status) > maxWidth {
			maxWidth = len(status)
		}
	}
	for i, status := range m.files {
		style := Modified
		if strings.HasPrefix(status, "M") {
			style = Modified
			if m.selected == i {
				style = ModifiedSelected
			}
		} else if strings.HasPrefix(status, "D") {
			style = Deleted
			if m.selected == i {
				style = DeletedSelected
			}
		} else {
			style = Added
			if m.selected == i {
				style = AddedSelected
			}
		}
		w.WriteString(style.Width(maxWidth).Render(status))
		w.WriteString("\n")
	}
	return w.String()
}
