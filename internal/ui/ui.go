package ui

import (
	"strings"

	"jjui/internal/dag"
	"jjui/internal/ui/common"
	"jjui/internal/ui/revisions"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	revisions revisions.Model
	help      help.Model
	width     int
	height    int
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle("jjui"), m.revisions.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	m.revisions, cmd = m.revisions.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.height == 0 {
		return "loading"
	}

	var b strings.Builder
	b.WriteString(m.help.View(&m.revisions.Keymap))
	b.WriteString("\n")

	footer := b.String()
	footerHeight := lipgloss.Height(footer)
	result := lipgloss.Place(m.width, m.height-footerHeight, 0, 0, m.revisions.View())
	return lipgloss.JoinVertical(0, result, footer)
}

func New() tea.Model {
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	return Model{
		revisions: revisions.New([]dag.GraphRow{}),
		help:      h,
	}
}
