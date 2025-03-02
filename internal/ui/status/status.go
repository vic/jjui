package status

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
)

type Model struct {
	spinner spinner.Model
	mode    string
	command string
	running bool
	output  string
	error   error
	width   int
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return 1
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetHeight(int) {}

func (m *Model) SetMode(mode string) {
	m.mode = mode
}

var (
	normalStyle  = lipgloss.NewStyle()
	successStyle = lipgloss.NewStyle().Inherit(normalStyle).Foreground(common.Green)
	errorStyle   = lipgloss.NewStyle().Inherit(normalStyle).Foreground(common.Red)
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.CommandRunningMsg:
		m.command = string(msg)
		m.running = true
		return m, m.spinner.Tick
	case common.CommandCompletedMsg:
		m.running = false
		m.output = msg.Output
		m.error = msg.Err
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	s := normalStyle.Render(" ")
	if !m.running {
		if m.error != nil {
			s = errorStyle.Render("✗ ")
		} else if m.command != "" {
			s = successStyle.Render("✓ ")
		}
	} else {
		s = normalStyle.Render(m.spinner.View())
	}
	ret := normalStyle.Width(m.width - 2).SetString(m.command).Render()
	ret = lipgloss.JoinHorizontal(lipgloss.Left, m.mode, s, ret)
	if m.error != nil {
		ret += " " + errorStyle.Render(fmt.Sprintf("\n%v\n%s", m.error, m.output))
	}
	return ret
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		spinner: s,
		command: "",
		running: false,
		output:  "",
	}
}
