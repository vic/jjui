package revisions

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"jjui/internal/jj"
	"jjui/internal/ui/abandon"
	"jjui/internal/ui/bookmark"
	"jjui/internal/ui/common"
	"jjui/internal/ui/describe"
	"jjui/internal/ui/revisions/revset"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type viewRange struct {
	start int
	end   int
}

type Model struct {
	rows        []jj.GraphRow
	status      common.Status
	error       error
	op          common.Operation
	viewRange   *viewRange
	draggedRow  int
	cursor      int
	Width       int
	Height      int
	overlay     tea.Model
	revsetModel revset.Model
	Keymap      keymap
	common.Commands
}

func (m Model) selectedRevision() *jj.Commit {
	if m.cursor >= len(m.rows) {
		return nil
	}
	return m.rows[m.cursor].Commit
}

func (m Model) Init() tea.Cmd {
	return common.Refresh("@")
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
	case key.Matches(msg, layer.revset):
		m.revsetModel, _ = m.revsetModel.Update(revset.EditRevSetMsg{})
		return m, nil

	case key.Matches(msg, layer.new):
		return m, tea.Sequence(
			m.NewRevision(m.selectedRevision().ChangeId),
			common.Refresh("@"),
		)
	case key.Matches(msg, layer.edit):
		return m, tea.Sequence(
			m.Edit(m.selectedRevision().ChangeId),
			common.Refresh("@"),
		)
	case key.Matches(msg, layer.diffedit):
		return m, m.DiffEdit(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.abandon):
		m.overlay = abandon.New(m.selectedRevision().ChangeId)
		return m, m.overlay.Init()
	case key.Matches(msg, layer.split):
		currentRevision := m.selectedRevision().ChangeId
		return m, m.Split(currentRevision)
	case key.Matches(msg, layer.description):
		m.overlay = describe.New(m.selectedRevision().ChangeId, m.selectedRevision().Description, m.Width)
		m.op = common.EditDescriptionOperation
		return m, m.overlay.Init()
	case key.Matches(msg, layer.diff):
		return m, m.GetDiff(m.selectedRevision().ChangeId)
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
		m.op = common.RebaseRevisionOperation
		m.draggedRow = m.cursor
	case key.Matches(msg, layer.branch):
		m.op = common.RebaseBranchOperation
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
		m.Keymap.resetMode()
		if m.draggedRow == -1 {
			m.op = common.None
			break
		}
		rebaseOperation := m.op
		fromRevision := m.rows[m.draggedRow].Commit.ChangeIdShort
		toRevision := m.rows[m.cursor].Commit.ChangeIdShort
		m.op = common.None
		m.draggedRow = -1
		return m, tea.Sequence(
			m.Rebase(fromRevision, toRevision, rebaseOperation),
			common.Refresh(fromRevision),
		)
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
		return m, m.FetchBookmarks(m.selectedRevision().ChangeId)
	case key.Matches(msg, layer.set):
		m.overlay = bookmark.NewSetBookmark(m.selectedRevision().ChangeId)
		m.op = common.SetBookmarkOperation
		return m, m.overlay.Init()
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
		return m, tea.Sequence(
			m.GitFetch(),
			common.Refresh(m.selectedRevision().ChangeId),
		)
	case key.Matches(msg, layer.push):
		m.Keymap.resetMode()
		return m, tea.Sequence(
			m.GitPush(),
			common.Refresh(m.selectedRevision().ChangeId),
		)
	case key.Matches(msg, m.Keymap.cancel):
		m.Keymap.resetMode()
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case common.CloseViewMsg:
		m.Keymap.resetMode()
		m.overlay = nil
		m.op = common.None
		return m, nil
	case common.UpdateRevSetMsg:
		m.revsetModel.Value = string(msg)
		if selectedRevision := m.selectedRevision(); selectedRevision != nil {
			return m, common.Refresh(selectedRevision.ChangeId)
		}
		return m, common.Refresh("@")
	case common.RefreshMsg:
		return m, tea.Sequence(
			m.FetchRevisions(m.revsetModel.Value),
			common.SelectRevision(msg.SelectedRevision),
		)
	case common.AbandonMsg:
		return m, tea.Sequence(
			m.Abandon(string(msg)),
			common.Refresh("@"),
			common.Close,
		)
	case common.MoveBookmarkMsg:
		return m, tea.Sequence(
			m.MoveBookmark(msg.Revision, msg.Bookmark),
			common.Refresh(msg.Revision),
			common.Close,
		)
	case common.SetBookmarkMsg:
		return m, tea.Sequence(
			m.SetBookmark(msg.Revision, msg.Bookmark),
			common.Refresh(msg.Revision),
			common.Close,
		)
	case common.SetDescriptionMsg:
		return m, tea.Sequence(
			m.SetDescription(msg.Revision, msg.Description),
			common.Refresh(msg.Revision),
			common.Close,
		)
	}

	var cmd tea.Cmd
	if m.overlay != nil {
		m.overlay, cmd = m.overlay.Update(msg)
		return m, cmd
	}

	if m.revsetModel.Editing {
		m.revsetModel, cmd = m.revsetModel.Update(msg)
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
			return row.Commit.ChangeId == r
		})
		if idx != -1 {
			m.cursor = idx
		} else {
			m.cursor = 0
		}
	case common.UpdateRevisionsMsg:
		if msg != nil {
			m.rows = msg
			m.status = common.Ready
			m.cursor = 0
		}
	case common.UpdateRevisionsFailedMsg:
		if msg != nil {
			m.rows = nil
			m.status = common.Error
			m.error = msg
		}
	case common.UpdateBookmarksMsg:
		m.overlay = bookmark.New(m.selectedRevision().ChangeId, msg, m.Width)
		return m, m.overlay.Init()
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}
	return m, cmd
}

func (m Model) View() string {
	revset := m.revsetModel.View()
	content := ""

	if m.status == common.Loading {
		content = "loading"
	}

	if m.status == common.Error {
		content = fmt.Sprintf("error: %v", m.error)
	}

	if m.status == common.Ready {
		h := m.Height - lipgloss.Height(revset)

		if m.viewRange.end-m.viewRange.start > h {
			m.viewRange.end = m.viewRange.start + h
		}

		highlightedRevision := ""
		if m.cursor < len(m.rows) {
			highlightedRevision = m.rows[m.cursor].Commit.ChangeIdShort
		}

		nodeRenderer := SegmentedRenderer{
			Palette:             common.DefaultPalette,
			op:                  m.op,
			HighlightedRevision: highlightedRevision,
			Overlay:             m.overlay,
		}

		var w jj.GraphWriter
		selectedLineStart := -1
		selectedLineEnd := -1
		for i, row := range m.rows {
			if i == m.cursor {
				selectedLineStart = w.LineCount()
			}
			w.RenderRow(row, nodeRenderer)
			if i == m.cursor {
				selectedLineEnd = w.LineCount()
			}
			if selectedLineEnd > 0 && w.LineCount() > h && w.LineCount() > m.viewRange.end {
				break
			}
		}

		if selectedLineStart <= m.viewRange.start {
			m.viewRange.start = selectedLineStart
			m.viewRange.end = selectedLineStart + h
		} else if selectedLineEnd > m.viewRange.end {
			m.viewRange.end = selectedLineEnd
			m.viewRange.start = selectedLineEnd - h
		}

		content = w.String(m.viewRange.start, m.viewRange.end)
	}

	return lipgloss.JoinVertical(0, revset, content)
}

func New(jj jj.Commands) Model {
	v := viewRange{start: 0, end: 0}
	defaultRevSet, _ := jj.GetConfig("revsets.log")
	return Model{
		status:      common.Loading,
		Commands:    common.NewCommands(jj),
		rows:        nil,
		draggedRow:  -1,
		viewRange:   &v,
		op:          common.None,
		cursor:      0,
		Width:       20,
		revsetModel: revset.New(string(defaultRevSet)),
		Keymap:      newKeyMap(),
	}
}
