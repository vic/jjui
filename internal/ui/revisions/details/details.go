package details

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"jjui/internal/ui/common"
)

var (
	AddedStyle    = lipgloss.NewStyle().Foreground(common.Green)
	DeletedStyle  = lipgloss.NewStyle().Foreground(common.Red)
	ModifiedStyle = lipgloss.NewStyle().Foreground(common.Cyan)
)

type status uint8

var (
	Added    status = 0
	Deleted  status = 1
	Modified status = 2
)

type item struct {
	status   status
	name     string
	selected bool
}

func (f item) Title() string       { return fmt.Sprintf("%c %s", f.status, f.name) }
func (f item) Description() string { return "" }
func (f item) FilterValue() string { return f.name }

type itemDelegate struct{}

func (i itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(item)
	if !ok {
		return
	}
	var style lipgloss.Style
	switch item.status {
	case Added:
		style = AddedStyle
	case Deleted:
		style = DeletedStyle
	case Modified:
		style = ModifiedStyle
	}
	if index == m.Index() {
		style = style.Bold(true).Background(common.DarkBlack)
	}
	title := item.Title()
	if item.selected {
		title = "✓" + title
	} else {
		title = " " + title
	}

	fmt.Fprint(w, style.Render(title))
}

func (i itemDelegate) Height() int                         { return 1 }
func (i itemDelegate) Spacing() int                        { return 0 }
func (i itemDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

type Model struct {
	revision string
	files    list.Model
	common.Commands
}

func New(revision string, commands common.Commands) tea.Model {
	l := list.New(nil, itemDelegate{}, 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	return Model{
		revision: revision,
		files:    l,
		Commands: commands,
	}
}

func (m Model) Init() tea.Cmd {
	return m.Status(m.revision)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "h":
			return m, common.Close
		case "d":
			v := m.files.SelectedItem().(item).name
			return m, m.Commands.GetDiff(m.revision, v)
		case "r":
			selectedFiles := make([]string, 0)
			for _, f := range m.files.Items() {
				if f.(item).selected {
					selectedFiles = append(selectedFiles, f.(item).name)
				}
			}
			if len(selectedFiles) == 0 {
				return m, nil
			}
			return m, m.Commands.Restore(m.revision, selectedFiles)
		case " ", "m":
			item := m.files.SelectedItem().(item)
			item.selected = !item.selected
			return m, m.files.SetItem(m.files.Index(), item)
		default:
			var cmd tea.Cmd
			m.files, cmd = m.files.Update(msg)
			return m, cmd
		}
	case common.RefreshMsg:
		return m, m.Status(m.revision)
	case common.UpdateCommitStatusMsg:
		items := make([]list.Item, 0)
		for _, file := range msg {
			if file == "" {
				continue
			}
			var status status
			switch file[0] {
			case 'A':
				status = Added
			case 'D':
				status = Deleted
			case 'M':
				status = Modified
			}
			items = append(items, item{
				status: status,
				name:   file[2:],
			})
		}
		m.files.SetItems(items)
		m.files.SetHeight(min(10, len(items)))
	}
	return m, nil
}

func (m Model) View() string {
	return m.files.View()
}
