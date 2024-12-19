package bookmark

import (
	"fmt"
	"io"
	"strings"

	"jjui/internal/ui/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	revision string
	list     list.Model
	op       common.Operation
}

type item string

func (b item) Title() string       { return string(b) }
func (b item) Description() string { return "" }
func (b item) FilterValue() string { return string(b) }

type itemDelegate struct{}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	str := fmt.Sprintf("%d %s", index+1, i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, common.Close
		case "enter":
			bookmark := m.list.SelectedItem().(item)
			switch m.op {
			case common.MoveBookmarkOperation:
				return m, moveBookmark(m.revision, string(bookmark))
			case common.DeleteBookmarkOperation:
				return m, deleteBookmark(string(bookmark))
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func moveBookmark(revision string, bookmark string) tea.Cmd {
	return func() tea.Msg {
		return common.MoveBookmarkMsg{
			Revision: revision,
			Bookmark: bookmark,
		}
	}
}

func deleteBookmark(bookmark string) tea.Cmd {
	return func() tea.Msg {
		return common.DeleteBookmarkMsg{
			Bookmark: bookmark,
		}
	}
}
func New(revision string, bookmarks []string, op common.Operation, width int) tea.Model {
	var items []list.Item
	for _, bookmark := range bookmarks {
		item := item(bookmark)
		items = append(items, item)
	}

	l := list.New(items, itemDelegate{}, width, len(items)*2)
	l.SetFilteringEnabled(true)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	return Model{
		revision: revision,
		op:       op,
		list:     l,
	}
}
