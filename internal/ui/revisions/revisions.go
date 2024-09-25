package revisions

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"jjui/internal/dag"
	"jjui/internal/jj"
	"jjui/internal/ui/bookmark"
	"jjui/internal/ui/common"
	"jjui/internal/ui/describe"
	"jjui/internal/ui/diff"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	rows       []dag.GraphRow
	op         common.Operation
	draggedRow int
	cursor     int
	width      int
	height     int
	help       help.Model
	overlay    tea.Model
	diff       tea.Model
	output     string
	keymap     keymap
}

func (m Model) selectedRevision() *jj.Commit {
	return m.rows[m.cursor].Commit
}

func (m Model) Init() tea.Cmd {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return tea.Sequence(tea.SetWindowTitle("jjui"), common.FetchRevisions(dir), common.SelectRevision("@"))
}

func (m Model) handleBaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(baseLayer)
	switch {
	case key.Matches(msg, layer.quit):
		return m, tea.Quit
	case key.Matches(msg, layer.gitMode):
		m.keymap.gitMode()
		return m, nil
	case key.Matches(msg, layer.rebaseMode):
		m.keymap.rebaseMode()
		return m, nil
	case key.Matches(msg, layer.bookmarkMode):
		m.keymap.bookmarkMode()
		return m, nil
	case key.Matches(msg, layer.description):
		return m, common.ShowDescribe(m.selectedRevision())
	case key.Matches(msg, layer.diff):
		return m, common.GetDiff(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.new):
		return m, common.NewRevision(m.selectedRevision().ChangeId)
	case key.Matches(msg, m.keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.keymap.cancel):
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleRebaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(rebaseLayer)
	switch {
	case key.Matches(msg, layer.revision):
		m.op = common.RebaseRevision
		m.draggedRow = m.cursor
	case key.Matches(msg, layer.branch):
		m.keymap.current = 'm'
		m.op = common.RebaseBranch
		m.draggedRow = m.cursor
	case key.Matches(msg, m.keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.keymap.apply):
		m.op = common.None
		fromRevision := m.rows[m.draggedRow].Commit.ChangeIdShort
		toRevision := m.rows[m.cursor].Commit.ChangeIdShort
		m.draggedRow = -1
		m.keymap.current = ' '
		return m, common.RebaseCommand(fromRevision, toRevision, m.op)
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
		m.op = common.None
	}
	return m, nil
}

func (m Model) handleBookmarkKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(bookmarkLayer)
	switch {
	case key.Matches(msg, layer.set):
		m.keymap.resetMode()
		return m, common.FetchBookmarks
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
	}
	return m, nil
}

func (m Model) handleGitKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.keymap.bindings[m.keymap.current].(gitLayer)
	switch {
	case key.Matches(msg, layer.fetch):
		m.keymap.resetMode()
		return m, common.GitFetch()
	case key.Matches(msg, layer.push):
		m.keymap.resetMode()
		return m, common.GitPush()
	case key.Matches(msg, m.keymap.cancel):
		m.keymap.resetMode()
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok {
		m.overlay = nil
		m.diff = nil
		return m, nil
	}

	var cmd tea.Cmd
	if m.overlay != nil {
		m.overlay, cmd = m.overlay.Update(msg)
		return m, cmd
	}

	if m.diff != nil {
		m.diff, cmd = m.diff.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.keymap.current {
		case ' ':
			return m.handleBaseKeys(msg)
		case 'r':
			return m.handleRebaseKeys(msg)
		case 'b':
			return m.handleBookmarkKeys(msg)
		case 'g':
			return m.handleGitKeys(msg)
		}

	case common.SelectRevisionMsg:
		r := string(msg)
		idx := slices.IndexFunc(m.rows, func(row dag.GraphRow) bool {
			if r == "@" {
				return row.Commit.IsWorkingCopy
			}
			return row.Commit.ChangeIdShort == r
		})
		if idx != -1 {
			m.cursor = idx
		} else {
			m.cursor = 0
		}
	case common.UpdateRevisionsMsg:
		m.rows = msg
	case common.UpdateBookmarksMsg:
		m.overlay = bookmark.New(m.selectedRevision().ChangeId, msg, m.width)
		return m, m.overlay.Init()
	case common.ShowDescribeViewMsg:
		m.overlay = describe.New(msg.ChangeId, msg.Description, m.width)
		return m, m.overlay.Init()
	case common.ShowDiffMsg:
		m.diff = diff.New(string(msg), m.width, m.height)
		return m, m.diff.Init()
	case common.ShowOutputMsg:
		m.output = msg.Output
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m Model) View() string {
	if m.height == 0 {
		return "loading"
	}

	if m.diff != nil {
		return m.diff.View()
	}

	var b strings.Builder
	b.WriteString(m.help.View(&m.keymap))
	b.WriteString("\n")
	if m.op == common.RebaseBranch || m.op == common.RebaseRevision {
		command := "-r"
		if m.op == common.RebaseBranch {
			command = "-b"
		}
		b.WriteString("jj rebase " + command + " " + m.rows[m.draggedRow].Commit.ChangeIdShort + " -d " + m.rows[m.cursor].Commit.ChangeIdShort + "\n")
	}

	if m.output != "" {
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("%v\n", m.output))
	}

	footer := b.String()
	footerHeight := lipgloss.Height(footer)
	if m.overlay != nil {
		footerHeight += 4
	}
	itemsPerPage := (m.height - footerHeight - 6) / 2
	if itemsPerPage == 0 {
		return b.String()
	}
	currentPage := m.cursor / itemsPerPage
	viewStart := currentPage * itemsPerPage
	viewEnd := (currentPage + 1) * itemsPerPage

	var items strings.Builder
	for i := 0; i < len(m.rows); i++ {
		if i < viewStart {
			continue
		}
		if i > viewEnd {
			continue
		}
		row := &m.rows[i]
		switch m.op {
		case common.RebaseRevision, common.RebaseBranch:
			if i == m.cursor {
				indent := strings.Repeat("│ ", row.Level)
				items.WriteString(indent)
				items.WriteString(common.DropStyle.Render("<< here >>"))
				items.WriteString("\n")
			}
			SegmentedRenderer(&items, row, common.DefaultPalette, i == m.draggedRow,
				Indent{}, NodeGlyph{}, "  ", ChangeIdShort{}, ChangeIdRest{}, " ", Author{}, " ", Branches{}, ConflictMarker{}, "\n",
				Indent{}, Glyph{}, "  ", Description{}, "\n",
				Indent{}, ElidedRevisions{})
		case common.None:
			SegmentedRenderer(&items, row, common.DefaultPalette, i == m.cursor,
				Indent{}, NodeGlyph{}, "  ", ChangeIdShort{}, ChangeIdRest{}, " ", Author{}, " ", Branches{}, ConflictMarker{}, "\n",
				Indent{}, Glyph{}, "  ", Description{}, "\n",
				Indent{}, ElidedRevisions{})
			if m.overlay != nil && m.cursor == i {
				items.WriteString(m.overlay.View())
				items.WriteString("\n")
			}
		}
	}
	items.WriteString("\n")
	items.WriteString(footer)
	return items.String()
}

func New(rows []dag.GraphRow) tea.Model {
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	return Model{
		rows:       rows,
		draggedRow: -1,
		op:         common.None,
		cursor:     0,
		width:      20,
		keymap:     newKeyMap(),
		help:       h,
	}
}
