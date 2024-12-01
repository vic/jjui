package revisions

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"jjui/internal/jj"
	"jjui/internal/ui/abandon"
	"jjui/internal/ui/bookmark"
	"jjui/internal/ui/common"
	"jjui/internal/ui/describe"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type viewRange struct {
	start int
	end   int
}

type Model struct {
	rows       []jj.GraphRow
	op         common.Operation
	viewRange  *viewRange
	draggedRow int
	cursor     int
	width      int
	height     int
	overlay    tea.Model
	Keymap     keymap
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
	layer := m.Keymap.bindings[m.Keymap.current].(baseLayer)
	switch {
	case key.Matches(msg, m.Keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.Keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, layer.new):
		return m, common.NewRevision(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.edit):
		return m, common.Edit(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.diffedit):
		return m, common.DiffEdit(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.abandon):
		m.overlay = abandon.New(m.selectedRevision().ChangeId)
		return m, m.overlay.Init()
	case key.Matches(msg, layer.split):
		return m, common.Split(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.description):
		m.overlay = describe.New(m.selectedRevision().ChangeId, m.selectedRevision().Description, m.width)
		return m, m.overlay.Init()
	case key.Matches(msg, layer.diff):
		return m, common.GetDiff(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.gitMode):
		m.Keymap.gitMode()
		return m, nil
	case key.Matches(msg, layer.rebaseMode):
		m.Keymap.rebaseMode()
		return m, nil
	case key.Matches(msg, layer.bookmarkMode):
		m.Keymap.bookmarkMode()
		return m, nil
	case key.Matches(msg, layer.quit), key.Matches(msg, m.Keymap.cancel):
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleRebaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.Keymap.bindings[m.Keymap.current].(rebaseLayer)
	switch {
	case key.Matches(msg, layer.revision):
		m.op = common.RebaseRevision
		m.draggedRow = m.cursor
	case key.Matches(msg, layer.branch):
		m.op = common.RebaseBranch
		m.draggedRow = m.cursor
	case key.Matches(msg, m.Keymap.down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.Keymap.up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.Keymap.apply):
		rebaseOperation := m.op
		fromRevision := m.rows[m.draggedRow].Commit.ChangeIdShort
		toRevision := m.rows[m.cursor].Commit.ChangeIdShort
		m.op = common.None
		m.draggedRow = -1
		m.Keymap.current = ' '
		return m, common.Rebase(fromRevision, toRevision, rebaseOperation)
	case key.Matches(msg, m.Keymap.cancel):
		m.Keymap.resetMode()
		m.op = common.None
	}
	return m, nil
}

func (m Model) handleBookmarkKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.Keymap.bindings[m.Keymap.current].(bookmarkLayer)
	switch {
	case key.Matches(msg, layer.move):
		m.Keymap.resetMode()
		return m, common.FetchBookmarks(m.selectedRevision().ChangeId)
	case key.Matches(msg, m.Keymap.cancel):
		m.Keymap.resetMode()
	}
	return m, nil
}

func (m Model) handleGitKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.Keymap.bindings[m.Keymap.current].(gitLayer)
	switch {
	case key.Matches(msg, layer.fetch):
		m.Keymap.resetMode()
		return m, common.GitFetch()
	case key.Matches(msg, layer.push):
		m.Keymap.resetMode()
		return m, common.GitPush()
	case key.Matches(msg, m.Keymap.cancel):
		m.Keymap.resetMode()
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok {
		m.overlay = nil
		return m, nil
	}

	var cmd tea.Cmd
	if m.overlay != nil {
		m.overlay, cmd = m.overlay.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Keymap.current {
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
		idx := slices.IndexFunc(m.rows, func(row jj.GraphRow) bool {
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, cmd
}

func (m Model) View() string {
	space := m.height

	if m.cursor < m.viewRange.start && m.viewRange.start > 0 {
		m.viewRange.start--
	}
	if m.cursor > m.viewRange.end && m.cursor < len(m.rows) {
		m.viewRange.start++
	}

	var items strings.Builder
	for i := m.viewRange.start; i < len(m.rows); i++ {
		if space <= 4 {
			break
		}
		line := strings.Builder{}
		row := &m.rows[i]
		switch m.op {
		case common.RebaseRevision, common.RebaseBranch:
			if i == m.cursor {
				indent := strings.Repeat("│ ", row.Level)
				line.WriteString(indent)
				line.WriteString(common.DropStyle.Render("<< here >>"))
				line.WriteString("\n")
			}
		}
		SegmentedRenderer(&line, row, common.DefaultPalette, i == m.cursor,
			Separate(" ", NodeGlyph{}, ChangeId{}, Author{}, Timestamp{}, Branches{}, ConflictMarker{}), "\n",
			Separate(" ", Glyph{}, If(m.overlay == nil || i != m.cursor, If(row.Commit.Empty, Empty{}, " "), Description{}), If(m.overlay != nil && i == m.cursor, Overlay(m.overlay))), "\n",
			ElidedRevisions{})
		s := line.String()
		space -= lipgloss.Height(s) - 1
		m.viewRange.end = i
		items.WriteString(s)
	}
	return items.String()
}

func New(rows []jj.GraphRow) Model {
	v := viewRange{start: 0, end: 0}
	return Model{
		rows:       rows,
		draggedRow: -1,
		viewRange:  &v,
		op:         common.None,
		cursor:     0,
		width:      20,
		Keymap:     newKeyMap(),
	}
}
