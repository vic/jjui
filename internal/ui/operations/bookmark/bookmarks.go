package bookmark

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	apply  = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "delete bookmark"))
)

type Model struct {
	revision string
	list     list.Model
	op       operations.Operation
	common.UICommands
}

type item string

func (b item) Title() string       { return string(b) }
func (b item) Description() string { return "" }
func (b item) FilterValue() string { return string(b) }

type itemDelegate struct{}

var (
	itemStyle = common.DefaultPalette.BookmarksStyle.PaddingLeft(1).PaddingRight(1)
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	if cur, ok := listItem.(item); ok {
		style := itemStyle
		if index == m.Index() {
			style = style.Bold(true).Background(common.Blue).Foreground(common.Black)
		}
		fmt.Fprint(w, "  "+style.Render(cur.Title()))
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cancel):
			return m, common.Close
		case key.Matches(msg, apply):
			bookmark := m.list.SelectedItem().(item)
			switch m.op.(type) {
			case MoveBookmarkOperation:
				return m, m.MoveBookmark(m.revision, string(bookmark))
			case DeleteBookmarkOperation:
				return m, m.DeleteBookmark(m.revision, string(bookmark))
			}
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

func New(commands common.UICommands, revision string) tea.Model {
	l := list.New(nil, itemDelegate{}, 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(true)
	l.SetShowHelp(false)
	l.Title = "Move Bookmark"
	return Model{
		revision:   revision,
		op:         MoveBookmarkOperation{},
		list:       l,
		UICommands: commands,
	}
}

func NewDeleteBookmark(commands common.UICommands, revision string, bookmarks []string, width int) tea.Model {
	items := convertToItems(bookmarks)
	l := list.New(items, itemDelegate{}, 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowTitle(true)
	l.Title = "Delete Bookmark"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(common.White)
	return Model{
		revision:   revision,
		op:         DeleteBookmarkOperation{},
		list:       l,
		UICommands: commands,
	}
}
