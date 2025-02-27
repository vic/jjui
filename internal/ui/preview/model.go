package preview

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"time"
)

type refreshPreviewContentMsg struct {
	Tag int
}

type updatePreviewContentMsg struct {
	Content string
}

type Model struct {
	tag     int
	view    viewport.Model
	help    help.Model
	width   int
	height  int
	content string
	context common.AppContext
}

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{
		m.view.KeyMap.HalfPageUp,
		m.view.KeyMap.HalfPageDown,
	}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return m.height
}

func (m *Model) SetWidth(w int) {
	content := lipgloss.NewStyle().MaxWidth(w - 2).Render(m.content)
	m.view.SetContent(content)
	m.view.Width = w
	m.width = w
}

func (m *Model) SetHeight(h int) {
	m.view.Height = h
	m.height = h
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updatePreviewContentMsg:
		m.content = msg.Content
		content := lipgloss.NewStyle().MaxWidth(m.Width() - 4).Render(msg.Content)
		m.view.SetContent(content)
		m.view.GotoTop()
	case common.SelectionChangedMsg, common.RefreshMsg:
		m.tag++
		tag := m.tag
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return refreshPreviewContentMsg{Tag: tag}
		})
	case refreshPreviewContentMsg:
		if m.tag == msg.Tag {
			switch msg := m.context.SelectedItem().(type) {
			case common.SelectedFile:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Diff(msg.ChangeId, msg.File))
					return updatePreviewContentMsg{Content: string(output)}
				}
			case common.SelectedRevision:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Show(msg.ChangeId))
					return updatePreviewContentMsg{Content: string(output)}
				}
			}
		}
	default:
		var cmd tea.Cmd
		m.view, cmd = m.view.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() string {
	return m.view.View()
}

func New(context common.AppContext) Model {
	view := viewport.New(0, 0)
	view.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	view.KeyMap.Up = key.NewBinding(key.WithDisabled())
	view.KeyMap.Down = key.NewBinding(key.WithKeys("J", "ctrl+n"))
	view.KeyMap.HalfPageDown = key.NewBinding(key.WithKeys("ctrl+d"))
	view.KeyMap.HalfPageUp = key.NewBinding(key.WithKeys("ctrl+u"))
	//view.KeyMap.PageUp = key.NewBinding(key.WithKeys("ctrl+u", "K"))
	//view.KeyMap.PageDown = key.NewBinding(key.WithKeys("ctrl+f", "J"))
	return Model{
		context: context,
		view:    view,
		help:    help.New(),
	}
}
