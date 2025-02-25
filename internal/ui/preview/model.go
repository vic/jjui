package preview

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"time"
)

type RefreshPreviewContentMsg struct {
	Tag      int
	Revision string
}

type Model struct {
	tag      int
	view     viewport.Model
	Width    int
	Height   int
	commands common.UICommands
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.UpdatePreviewContentMsg:
		content := lipgloss.NewStyle().MaxWidth(m.Width - 2).Render(msg.Content)
		m.view.SetContent(content)
	case RefreshPreviewContentMsg:
		if m.tag == msg.Tag {
			return m, m.commands.Show(msg.Revision)
		}
	case common.UpdatePreviewChangeIdMsg:
		m.tag++
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return RefreshPreviewContentMsg{Tag: m.tag, Revision: msg.ChangeId}
		})
	default:
		var cmd tea.Cmd
		m.view, cmd = m.view.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	m.view.Width = m.Width - 2
	m.view.Height = m.Height - 2
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(m.view.View())
}

func New(commands common.UICommands) Model {
	view := viewport.New(0, 0)
	view.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	view.KeyMap.Up = key.NewBinding(key.WithDisabled())
	view.KeyMap.Down = key.NewBinding(key.WithKeys("J", "ctrl+n"))
	view.KeyMap.HalfPageDown = key.NewBinding(key.WithKeys("ctrl+d"))
	view.KeyMap.HalfPageUp = key.NewBinding(key.WithKeys("ctrl+u"))
	//view.KeyMap.PageUp = key.NewBinding(key.WithKeys("ctrl+u", "K"))
	//view.KeyMap.PageDown = key.NewBinding(key.WithKeys("ctrl+f", "J"))
	return Model{
		commands: commands,
		view:     view,
	}
}
