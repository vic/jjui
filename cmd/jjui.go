package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/jj"
	"os"
	"strings"
)

type mode int

const (
	normalMode mode = iota
	moveMode
)

var highlightColor = lipgloss.Color("#282a36")
var commitShortStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#bd93f9"))

var commitIdRestStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#6272a4"))

var authorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ffb86c"))

var normal = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#f8f8f2"))

type model struct {
	items              []jj.Commit
	mode               mode
	draggedCommitIndex int
	cursor             int
	width              int
}

func fetchLog(location string) tea.Cmd {
	return func() tea.Msg {
		commits := jj.GetCommits(location)
		return logCommand(jj.BuildCommitTree(commits))
	}
}

func rebaseCommand(from, to string) tea.Cmd {
	if err := jj.RebaseCommand(from, to); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return fetchLog(os.Getenv("PWD"))
}

type logCommand []jj.Commit

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
				if m.cursor == m.draggedCommitIndex-1 || m.cursor == m.draggedCommitIndex {
					m.cursor++
				}
				if m.cursor < len(m.items) {
					m.cursor++
				}
			} else if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.mode == moveMode {
				if m.cursor == m.draggedCommitIndex+1 || m.cursor == m.draggedCommitIndex {
					m.cursor--
				}
				if m.cursor > 0 {
					m.cursor--
				}
			} else if m.cursor > 0 {
				m.cursor--
			}
		case "esc":
			m.draggedCommitIndex = -1
			m.mode = normalMode
		case " ":
			if m.mode == normalMode {
				m.mode = moveMode
				m.draggedCommitIndex = m.cursor
			} else {
				m.mode = normalMode
				m.draggedCommitIndex = -1
			}
		case "enter":
			if m.mode == moveMode {
				m.mode = normalMode
				fromRevision := m.items[m.draggedCommitIndex].ChangeIdShort
				toRevision := m.items[m.cursor].ChangeIdShort
				m.draggedCommitIndex = -1
				return m, rebaseCommand(fromRevision, toRevision)
			}
		default:
			return m, nil
		}
	case logCommand:
		commits := []jj.Commit(msg)
		m.items = commits
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m model) View() string {
	items := strings.Builder{}
	for i := 0; i < len(m.items); i++ {
		commit := &m.items[i]
		switch m.mode {
		case moveMode:
			if i == m.cursor {
				draggedCommit := &m.items[m.draggedCommitIndex]
				items.WriteString(m.viewCommit(draggedCommit, i == m.cursor, commit.Level()))
			}
			if i != m.draggedCommitIndex {
				items.WriteString(m.viewCommit(commit, false, commit.Level()))
			}
		case normalMode:
			items.WriteString(m.viewCommit(commit, i == m.cursor, commit.Level()))
		}
		if len(commit.Parents) == 0 && i < len(m.items)-1 {
			items.WriteString(commitIdRestStyle.Render(" ~ (elided revisions)"))
			items.WriteString("\n")
		}
	}
	if m.cursor == len(m.items) && m.mode == moveMode {
		items.WriteString(m.viewCommit(&m.items[m.draggedCommitIndex], true, m.items[m.draggedCommitIndex].Level()))
	}
	bottom := fmt.Sprintf("use j,k keys to move up and down: cursor:%v dragged:%d\n", m.cursor, m.draggedCommitIndex)
	if m.mode == moveMode {
		if m.cursor == len(m.items) {
			bottom += "jj rebase -r " + m.items[m.draggedCommitIndex].ChangeIdShort + " --insert-before " + m.items[len(m.items)-1].ChangeIdShort + "\n"
		} else {
			bottom += "jj rebase -r " + m.items[m.draggedCommitIndex].ChangeIdShort + " -d " + m.items[m.cursor].ChangeIdShort + "\n"
		}
	}
	items.WriteString(bottom)
	return items.String()
}

func (m model) viewCommit(commit *jj.Commit, highlighted bool, level int) string {
	changeIdRemaining := strings.TrimPrefix(commit.ChangeId, commit.ChangeIdShort)
	builder := strings.Builder{}
	for j := 0; j < level; j++ {
		builder.WriteString(normal.Render(" │ "))
	}

	if commit.IsWorkingCopy {
		builder.WriteString(normal.Render(" @ "))
	} else {
		builder.WriteString(normal.Render(" o "))
	}

	if highlighted {
		builder.WriteString(commitShortStyle.Background(highlightColor).Render(commit.ChangeIdShort))
		builder.WriteString(commitIdRestStyle.Background(highlightColor).Render(changeIdRemaining + " "))
		builder.WriteString(authorStyle.Background(highlightColor).Render(commit.Author) + "\n")
		builder.WriteString(strings.Repeat(" │ ", level+1))
		if commit.Description == "" {
			builder.WriteString(normal.Background(highlightColor).Bold(true).Foreground(lipgloss.Color("#50fa7b")).Width(m.width).Render("(no description)"))
		} else {
			builder.WriteString(normal.Background(highlightColor).Width(m.width).Render(commit.Description))
		}
	} else {
		builder.WriteString(commitShortStyle.Render(commit.ChangeIdShort))
		builder.WriteString(commitIdRestStyle.Render(changeIdRemaining + " "))
		builder.WriteString(authorStyle.Render(commit.Author) + "\n")
		builder.WriteString(strings.Repeat(" │ ", level+1))
		if commit.Description == "" {
			builder.WriteString(normal.Bold(true).Foreground(lipgloss.Color("#50fa7b")).Render("(no description)"))
		} else {
			builder.WriteString(normal.Render(commit.Description))
		}
	}
	builder.WriteString("\n")
	return builder.String()
}

func initialModel() model {
	return model{
		items:              []jj.Commit{},
		draggedCommitIndex: -1,
		mode:               normalMode,
		cursor:             0,
		width:              20,
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
