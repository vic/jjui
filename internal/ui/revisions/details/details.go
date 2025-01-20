package details

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"jjui/internal/ui/common"
	"strings"
)

var (
	Added    = lipgloss.NewStyle().Foreground(common.Green)
	Deleted  = lipgloss.NewStyle().Foreground(common.Red)
	Modified = lipgloss.NewStyle().Foreground(common.Cyan)
)

type item string

func (f item) Title() string       { return string(f) }
func (f item) Description() string { return "" }
func (f item) FilterValue() string { return string(f) }

type itemDelegate struct{}

func (i itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	f, ok := listItem.(item)
	if !ok {
		return
	}
	file := string(f)

	var style = Modified
	if strings.HasPrefix(file, "A") {
		style = Added
	} else if strings.HasPrefix(file, "D") {
		style = Deleted
	}
	if index == m.Index() {
		style = style.Bold(true).Background(common.DarkBlack)
	}
	fmt.Fprint(w, style.Render(file))
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
			v := m.files.SelectedItem().FilterValue()
			return m, m.Commands.GetDiff(m.revision, v[2:])
		default:
			var cmd tea.Cmd
			m.files, cmd = m.files.Update(msg)
			return m, cmd
		}
	case common.UpdateCommitStatusMsg:
		items := make([]list.Item, len(msg))
		for i, status := range msg {
			items[i] = item(status)
		}
		m.files.SetItems(items)
		m.files.SetHeight(min(10, len(items)))
	}
	return m, nil
}

func (m Model) View() string {
	return m.files.View()
}
