package confirmation

import (
	"github.com/charmbracelet/bubbles/key"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/idursun/jjui/internal/ui/common"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	right  = key.NewBinding(key.WithKeys("right", "l"))
	left   = key.NewBinding(key.WithKeys("left", "h"))
	enter  = key.NewBinding(key.WithKeys("enter"))
	cancel = key.NewBinding(key.WithKeys("esc"))
)

type CloseMsg struct{}

type option struct {
	label      string
	cmd        tea.Cmd
	keyBinding key.Binding
}

type Model struct {
	message  string
	options  []option
	selected int
}

var (
	textStyle   = lipgloss.NewStyle().Bold(true).Foreground(common.Magenta)
	normalStyle = lipgloss.NewStyle().
			Foreground(common.White).
			PaddingLeft(2).
			PaddingRight(2)
)

var selectedStyle = lipgloss.NewStyle().
	Foreground(common.DarkWhite).
	Background(common.Blue).
	PaddingLeft(2).
	PaddingRight(2)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, left):
			if m.selected > 0 {
				m.selected--
			}
		case key.Matches(msg, right):
			if m.selected < len(m.options) {
				m.selected++
			}
		case key.Matches(msg, enter):
			selectedOption := m.options[m.selected]
			return m, selectedOption.cmd
		default:
			for _, option := range m.options {
				if key.Matches(msg, option.keyBinding) {
					return m, option.cmd
				}
			}
		}
	}
	return m, nil
}

func (m *Model) View() string {
	w := strings.Builder{}
	w.WriteString(textStyle.Render(m.message))
	for i, option := range m.options {
		w.WriteString(" ")
		if i == m.selected {
			w.WriteString(selectedStyle.Render(option.label))
		} else {
			w.WriteString(normalStyle.Render(option.label))
		}
	}
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1, 0, 1).Render(w.String())
}

func (m *Model) AddOption(label string, cmd tea.Cmd, keyBinding key.Binding) {
	m.options = append(m.options, option{label, cmd, keyBinding})
}

func New(message string) Model {
	return Model{
		message:  message,
		options:  []option{},
		selected: 0,
	}
}

func Close() tea.Msg {
	return CloseMsg{}
}
