package revisions

import (
	"fmt"
	"jjui/internal/jj"
	"jjui/internal/ui/abandon"
	"jjui/internal/ui/bookmark"
	"jjui/internal/ui/common"
	"jjui/internal/ui/describe"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type viewRange struct {
	start int
	end   int
}

type Model struct {
	dag        *jj.Dag
	rows       []jj.TreeRow
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
	return &m.rows[m.cursor].Commit
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
		m.op = common.EditDescription
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
		m.op = common.None
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
		idx := slices.IndexFunc(m.rows, func(row jj.TreeRow) bool {
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
		if msg != nil {
			m.dag = msg
			m.rows = (*msg).GetTreeRows()
		}
	case common.UpdateBookmarksMsg:
		m.overlay = bookmark.New(m.selectedRevision().ChangeId, msg, m.width)
		return m, m.overlay.Init()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - 3
	}
	return m, cmd
}

func (m Model) View() string {
	if len(m.rows) == 0 {
		return "loading"
	}

	nodeRenderer := SegmentedRenderer{
		Palette:             common.DefaultPalette,
		op:                  m.op,
		HighlightedRevision: m.rows[m.cursor].Commit.ChangeIdShort,
		Overlay:             m.overlay,
	}

	var w jj.LineTrackingWriter
	selectedLineStart := -1
	selectedLineEnd := -1
	for i, row := range m.rows {
		if i == m.cursor {
			selectedLineStart = w.LineCount()
		}
		jj.RenderRow(&w, row, nodeRenderer)
		if i == m.cursor {
			selectedLineEnd = w.LineCount()
		}
	}

	if selectedLineStart <= m.viewRange.start {
		m.viewRange.start = selectedLineStart
		m.viewRange.end = selectedLineStart + m.height
	} else if selectedLineEnd > m.viewRange.end {
		m.viewRange.end = selectedLineEnd
		m.viewRange.start = selectedLineEnd - m.height
	}
	return w.String(m.viewRange.start, m.viewRange.end)
}

func New(dag *jj.Dag) Model {
	v := viewRange{start: 0, end: 0}
	var rows []jj.TreeRow
	if dag != nil {
		rows = dag.GetTreeRows()
	}
	return Model{
		dag:        dag,
		rows:       rows,
		draggedRow: -1,
		viewRange:  &v,
		op:         common.None,
		cursor:     0,
		width:      20,
		Keymap:     newKeyMap(),
	}
}
