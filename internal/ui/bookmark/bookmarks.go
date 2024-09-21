package bookmark

import (
	"fmt"
	"io"
	"strings"

	"jjui/internal/jj"
	"jjui/internal/ui/common"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	revision string
	list     list.Model
}

type item string

func (b item) Title() string       { return string(b) }
func (b item) Description() string { return "" }
func (b item) FilterValue() string { return string(b) }

type itemDelegate struct{}

var (
	itemSyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemSyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
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
	fn := itemSyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemSyle.Render("> " + strings.Join(s, " "))
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
			return m, tea.Batch(common.Close, common.SetBookmark(m.revision, string(bookmark)))
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

func New(revision string, bookmarks []jj.Bookmark, width int) tea.Model {
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
		list:     l,
	}
}
