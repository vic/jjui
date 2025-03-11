package bookmark

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
)

var (
	cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

type Model struct {
	revision string
	list     list.Model
	op       operations.Operation
	context  context.AppContext
}

type item string

func (b item) Title() string       { return string(b) }
func (b item) Description() string { return "" }
func (b item) FilterValue() string { return string(b) }

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cancel):
			return m, common.Close
		}
	case common.UpdateBookmarksMsg:
		items := convertToItems(msg.Bookmarks)
		m.list.SetItems(items)
		m.list.SetHeight(min(10, len(items)+2))
		return m, nil
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(min(10, len(m.list.Items())+2))
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func convertToItems(bookmarks []string) []list.Item {
	var items []list.Item
	for _, bookmark := range bookmarks {
		item := item(bookmark)
		items = append(items, item)
	}
	return items
}

func (m Model) View() string {
	return m.list.View()
}
