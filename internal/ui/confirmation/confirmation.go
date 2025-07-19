package confirmation

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
)

var (
	right = key.NewBinding(key.WithKeys("right", "l"))
	left  = key.NewBinding(key.WithKeys("left", "h"))
	enter = key.NewBinding(key.WithKeys("enter"))
)

type CloseMsg struct{}

type option struct {
	label      string
	cmd        tea.Cmd
	keyBinding key.Binding
}

type Styles struct {
	Border   lipgloss.Style
	Selected lipgloss.Style
	Dimmed   lipgloss.Style
	Text     lipgloss.Style
}

type Model struct {
	options  []option
	selected int
	Styles   Styles
	messages []string
}

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
			if m.selected < len(m.options)-1 {
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
	for i, message := range m.messages {
		w.WriteString(m.Styles.Text.Render(message))
		if i < len(m.messages)-1 {
			w.WriteString(m.Styles.Text.Render("\n"))
		}
	}
	for i, option := range m.options {
		if i == m.selected {
			w.WriteString(m.Styles.Selected.Render(option.label))
		} else {
			w.WriteString(m.Styles.Dimmed.Render(option.label))
		}
	}
	content := w.String()
	width, height := lipgloss.Size(content)
	content = lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content, lipgloss.WithWhitespaceBackground(m.Styles.Text.GetBackground()))
	return m.Styles.Border.Render(content)
}

func (m *Model) AddOption(label string, cmd tea.Cmd, keyBinding key.Binding) {
	m.options = append(m.options, option{label, cmd, keyBinding})
}

func New(messages ...string) Model {
	styles := Styles{
		Border:   common.DefaultPalette.GetBorder("confirmation border", lipgloss.RoundedBorder()),
		Text:     common.DefaultPalette.Get("confirmation text").PaddingRight(1),
		Selected: common.DefaultPalette.Get("confirmation selected").PaddingLeft(2).PaddingRight(2),
		Dimmed:   common.DefaultPalette.Get("confirmation dimmed").PaddingLeft(2).PaddingRight(2),
	}
	return Model{
		messages: messages,
		options:  []option{},
		selected: 0,
		Styles:   styles,
	}
}

func Close() tea.Msg {
	return CloseMsg{}
}
