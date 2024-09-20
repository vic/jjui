package revisions

import (
	"fmt"
	"os"
	"strings"

	"jjui/internal/dag"
	"jjui/internal/ui/common"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	normalMode mode = iota
	moveMode
)

type Model struct {
	rows       []dag.GraphRow
	mode       mode
	draggedRow int
	cursor     int
	width      int
	help       help.Model
	keymap     keymap
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "rebase start")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "rebase apply")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit/rebase cancel")),
		key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "update description")),
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func (m Model) Init() tea.Cmd {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return common.FetchRevisions(dir)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "d":
			return m, common.DoShowDescribe(m.rows[m.cursor].Commit)
		case "down", "j":
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "esc":
			m.draggedRow = -1
			m.mode = normalMode
		case " ":
			if m.mode == normalMode {
				m.mode = moveMode
				m.draggedRow = m.cursor
			} else {
				m.mode = normalMode
				m.draggedRow = -1
			}
		case "enter":
			if m.mode == moveMode {
				m.mode = normalMode
				fromRevision := m.rows[m.draggedRow].Commit.ChangeIdShort
				toRevision := m.rows[m.cursor].Commit.ChangeIdShort
				m.draggedRow = -1
				return m, common.RebaseCommand(fromRevision, toRevision)
			}
		default:
			return m, nil
		}
	case common.UpdateRevisions:
		rows := []dag.GraphRow(msg)
		m.rows = rows
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m Model) View() string {
	var items strings.Builder
	for i := 0; i < len(m.rows); i++ {
		row := &m.rows[i]
		switch m.mode {
		case moveMode:
			if i == m.cursor {
				indent := strings.Repeat("│ ", row.Level)
				items.WriteString(indent)
				items.WriteString(dag.DropStyle.Render("<< here >>"))
				items.WriteString("\n")
			}
			dag.DefaultRenderer(&items, row, dag.DefaultPalette, i == m.draggedRow)
		case normalMode:
			dag.DefaultRenderer(&items, row, dag.DefaultPalette, i == m.cursor)
		}
	}
	items.WriteString("\n")
	items.WriteString(m.help.View(m.keymap))
	items.WriteString("\n")
	if m.mode == moveMode {
		if m.cursor == len(m.rows) {
			items.WriteString("jj rebase -r " + m.rows[m.draggedRow].Commit.ChangeIdShort + " --insert-before " + m.rows[len(m.rows)-1].Commit.ChangeIdShort + "\n")
		} else {
			items.WriteString("jj rebase -r " + m.rows[m.draggedRow].Commit.ChangeIdShort + " -d " + m.rows[m.cursor].Commit.ChangeIdShort + "\n")
		}
	}
	return items.String()
}

func New() tea.Model {
	help := help.New()
	help.Styles.ShortKey = dag.DefaultPalette.CommitShortStyle
	help.Styles.ShortDesc = dag.DefaultPalette.CommitIdRestStyle
	return Model{
		rows:       []dag.GraphRow{},
		draggedRow: -1,
		mode:       normalMode,
		cursor:     0,
		width:      20,
		keymap:     keymap{},
		help:       help,
	}
}
