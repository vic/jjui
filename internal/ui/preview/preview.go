package preview

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
)

type Model struct {
	tag     int
	focused bool
	view    viewport.Model
	help    help.Model
	width   int
	height  int
	content string
	context context.AppContext
	keyMap  config.KeyMappings[key.Binding]
}

const DebounceTime = 200 * time.Millisecond

var tab = key.NewBinding(key.WithKeys("tab"))

type refreshPreviewContentMsg struct {
	Tag int
}

type updatePreviewContentMsg struct {
	Content string
}
type focusMsg struct{}

func (m *Model) IsFocused() bool {
	return m.focused
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

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updatePreviewContentMsg:
		m.content = msg.Content
		content := lipgloss.NewStyle().MaxWidth(m.Width() - 4).Render(msg.Content)
		m.view.SetContent(content)
		m.view.GotoTop()
	case common.SelectionChangedMsg, common.RefreshMsg:
		m.tag++
		tag := m.tag
		return m, tea.Tick(DebounceTime, func(t time.Time) tea.Msg {
			return refreshPreviewContentMsg{Tag: tag}
		})
	case focusMsg:
		m.focused = true
		m.view.KeyMap = viewport.DefaultKeyMap()
		return m, nil
	case refreshPreviewContentMsg:
		if m.tag == msg.Tag {
			switch msg := m.context.SelectedItem().(type) {
			case context.SelectedFile:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Diff(msg.ChangeId, msg.File))
					return updatePreviewContentMsg{Content: string(output)}
				}
			case context.SelectedRevision:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.Show(msg.ChangeId))
					return updatePreviewContentMsg{Content: string(output)}
				}
			case context.SelectedOperation:
				return m, func() tea.Msg {
					output, _ := m.context.RunCommandImmediate(jj.OpShow(msg.OperationId))
					return updatePreviewContentMsg{Content: string(output)}
				}
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Cancel), key.Matches(msg, tab):
			m.focused = false
			m.view.KeyMap = unfocusedKeyMap()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.view.View()
}

func Focus() tea.Msg {
	return focusMsg{}
}

func unfocusedKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		PageDown:     key.NewBinding(key.WithDisabled()),
		PageUp:       key.NewBinding(key.WithDisabled()),
		HalfPageUp:   key.NewBinding(key.WithKeys("ctrl+u")),
		HalfPageDown: key.NewBinding(key.WithKeys("ctrl+d")),
		Up:           key.NewBinding(key.WithKeys("ctrl+p")),
		Down:         key.NewBinding(key.WithKeys("ctrl+n")),
	}
}

func New(context context.AppContext) Model {
	view := viewport.New(0, 0)
	view.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	view.KeyMap = unfocusedKeyMap()
	return Model{
		context: context,
		keyMap:  context.KeyMap(),
		view:    view,
		help:    help.New(),
	}
}
