package ui

import (
	"jjui/internal/jj"
	"strings"

	"jjui/internal/ui/common"
	"jjui/internal/ui/diff"
	"jjui/internal/ui/revisions"
	"jjui/internal/ui/status"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	revisions revisions.Model
	diff      tea.Model
	help      help.Model
	status    status.Model
	output    string
	width     int
	height    int
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle("jjui"), m.revisions.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok && m.diff != nil {
		m.diff = nil
		return m, nil
	}

	var cmd tea.Cmd
	if m.diff != nil {
		m.diff, cmd = m.diff.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case common.ShowDiffMsg:
		m.diff = diff.New(string(msg), m.width, m.height)
		return m, m.diff.Init()
	case common.CommandCompletedMsg:
		m.output = msg.Output
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	m.revisions, cmd = m.revisions.Update(msg)
	var statusCmd tea.Cmd
	m.status, statusCmd = m.status.Update(msg)
	return m, tea.Batch(cmd, statusCmd)
}

func (m Model) View() string {
	if m.height == 0 {
		return "loading"
	}

	if m.diff != nil {
		return m.diff.View()
	}

	var b strings.Builder
	b.WriteString(m.help.View(&m.revisions.Keymap))
	b.WriteString("\n")
	b.WriteString(m.status.View())

	footer := b.String()
	footerHeight := lipgloss.Height(footer)
	result := lipgloss.Place(m.width, m.height-footerHeight, 0, 0, m.revisions.View())
	return lipgloss.JoinVertical(0, result, footer)
}

func New() tea.Model {
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	d := jj.NewDag()
	return Model{
		revisions: revisions.New(&d),
		help:      h,
		status:    status.New(),
	}
}
