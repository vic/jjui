package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/dag"
	"jjui/internal/jj"
	"os"
	"strings"
)

type mode int

const (
	normalMode mode = iota
	moveMode
)

type model struct {
	rows       []dag.GraphRow
	mode       mode
	draggedRow int
	cursor     int
	width      int
}

type logCommand []dag.GraphRow

func fetchLog(location string) tea.Cmd {
	return func() tea.Msg {
		commits := jj.GetCommits(location)
		root := dag.Build(commits)
		rows := dag.BuildGraphRows(root)
		return logCommand(rows)
	}
}

func rebaseCommand(from, to string) tea.Cmd {
	if err := jj.RebaseCommand(from, to); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return fetchLog(os.Getenv("PWD"))
}
func (m model) Init() tea.Cmd {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return fetchLog(dir)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "down", "j":
			if m.mode == moveMode {
				//skip over dragged commit
				if m.cursor == m.draggedRow-1 || m.cursor == m.draggedRow {
					m.cursor++
				}
				if m.cursor < len(m.rows) {
					m.cursor++
				}
			} else if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.mode == moveMode {
				if m.cursor == m.draggedRow+1 || m.cursor == m.draggedRow {
					m.cursor--
				}
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.cursor > 0 {
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
				return m, rebaseCommand(fromRevision, toRevision)
			}
		default:
			return m, nil
		}
	case logCommand:
		rows := []dag.GraphRow(msg)
		m.rows = rows
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	var items strings.Builder
	for i := 0; i < len(m.rows); i++ {
		row := &m.rows[i]
		switch m.mode {
		case moveMode:
			if i == m.cursor {
				draggedRow := &m.rows[m.draggedRow]
				dag.DefaultRenderer(&items, draggedRow, dag.HighlightedPalette)
			}
			if i != m.draggedRow {
				dag.DefaultRenderer(&items, row, dag.DefaultPalette)
			}
		case normalMode:
			palette := dag.DefaultPalette
			if i == m.cursor {
				palette = dag.HighlightedPalette
			}
			dag.DefaultRenderer(&items, row, palette)
		}
	}
	if m.cursor == len(m.rows) && m.mode == moveMode {
		//TODO: should be rendered at a different level
		dag.DefaultRenderer(&items, &m.rows[m.draggedRow], dag.HighlightedPalette)
	}
    items.WriteString(fmt.Sprintf("use j,k keys to move up and down: cursor:%v dragged:%d\n", m.cursor, m.draggedRow))
	if m.mode == moveMode {
		if m.cursor == len(m.rows) {
			items.WriteString("jj rebase -r " + m.rows[m.draggedRow].Commit.ChangeIdShort + " --insert-before " + m.rows[len(m.rows)-1].Commit.ChangeIdShort + "\n")
		} else {
			items.WriteString("jj rebase -r " + m.rows[m.draggedRow].Commit.ChangeIdShort + " -d " + m.rows[m.cursor].Commit.ChangeIdShort + "\n")
		}
	}
	return items.String()
}

func initialModel() model {
	return model{
		rows:       []dag.GraphRow{},
		draggedRow: -1,
		mode:       normalMode,
		cursor:     0,
		width:      20,
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
