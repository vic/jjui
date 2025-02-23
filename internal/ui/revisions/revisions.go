package revisions

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/idursun/jjui/internal/ui/keymap"
	"github.com/idursun/jjui/internal/ui/operations/git"
	"github.com/idursun/jjui/internal/ui/operations/rebase"
	"slices"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/operations"
	"github.com/idursun/jjui/internal/ui/operations/abandon"
	"github.com/idursun/jjui/internal/ui/operations/bookmark"
	"github.com/idursun/jjui/internal/ui/operations/describe"
	"github.com/idursun/jjui/internal/ui/operations/details"
	"github.com/idursun/jjui/internal/ui/operations/squash"
	"github.com/idursun/jjui/internal/ui/revisions/revset"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewRange struct {
	start int
	end   int
}

type Model struct {
	rows         []jj.GraphRow
	status       common.Status
	error        error
	op           operations.Operation
	viewRange    *viewRange
	cursor       int
	Width        int
	Height       int
	revsetModel  revset.Model
	confirmation *confirmation.Model
	Keymap       keymap.Keymap
	common.UICommands
}

func (m Model) selectedRevision() *jj.Commit {
	if m.cursor >= len(m.rows) {
		return nil
	}
	return m.rows[m.cursor].Commit
}

func (m Model) GetKeyMap() help.KeyMap {
	if op, ok := m.op.(help.KeyMap); ok {
		return op
	}
	return &m.Keymap
}

func (m Model) Init() tea.Cmd {
	return common.Refresh("@")
}

func (m Model) handleBaseKeys(msg tea.KeyMsg) (Model, tea.Cmd) {
	layer := m.Keymap.Bindings[m.Keymap.Current].(keymap.BaseLayer)
	switch {
	case key.Matches(msg, m.Keymap.Up):
		if m.cursor > 0 {
			m.cursor--
		}
	case key.Matches(msg, m.Keymap.Down):
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
	case key.Matches(msg, m.Keymap.Cancel):
		m.Keymap.ResetMode()
		m.op = &operations.Noop{}
	case key.Matches(msg, m.Keymap.Details):
		m.Keymap.DetailsMode()
		var cmd tea.Cmd
		m.op, cmd = details.NewOperation(m.UICommands, m.selectedRevision())
		return m, cmd
	case key.Matches(msg, layer.Revset):
		m.revsetModel, _ = m.revsetModel.Update(revset.EditRevSetMsg{})
		return m, nil
	case key.Matches(msg, layer.Undo):
		model := confirmation.New("Are you sure you want to undo last change?")
		model.AddOption("Yes", tea.Batch(m.Undo(), confirmation.Close))
		model.AddOption("No", confirmation.Close)
		m.confirmation = &model
		return m, m.confirmation.Init()
	case key.Matches(msg, layer.New):
		return m, m.NewRevision(m.selectedRevision().GetChangeId())
	case key.Matches(msg, layer.Edit):
		return m, m.Edit(m.selectedRevision().GetChangeId())
	case key.Matches(msg, layer.Diffedit):
		return m, m.DiffEdit(m.selectedRevision().GetChangeId())
	case key.Matches(msg, layer.Abandon):
		var cmd tea.Cmd
		m.op, cmd = abandon.NewOperation(m.UICommands, m.selectedRevision())
		return m, cmd
	case key.Matches(msg, layer.Split):
		currentRevision := m.selectedRevision().GetChangeId()
		return m, m.Split(currentRevision, []string{})
	case key.Matches(msg, layer.Description):
		var cmd tea.Cmd
		m.op, cmd = describe.NewOperation(m.UICommands, m.selectedRevision(), m.Width)
		return m, cmd
	case key.Matches(msg, layer.Diff):
		return m, m.GetDiff(m.selectedRevision().GetChangeId(), "")
	case key.Matches(msg, layer.Refresh):
		return m, common.Refresh(m.selectedRevision().GetChangeId())
	case key.Matches(msg, layer.GitMode):
		m.op = git.NewOperation(m.UICommands)
		return m, nil
	case key.Matches(msg, layer.SquashMode):
		m.op = squash.NewOperation(m.UICommands, m.selectedRevision().ChangeIdShort)
		if m.cursor < len(m.rows)-1 {
			m.cursor++
		}
		return m, nil
	case key.Matches(msg, layer.RebaseMode):
		m.op = rebase.NewOperation(m.UICommands, m.selectedRevision().ChangeIdShort, rebase.SourceRevision, rebase.TargetDestination)
		return m, nil
	case key.Matches(msg, layer.BookmarkMode):
		m.op = bookmark.NewChooseBookmarkOperation(m.UICommands)
		return m, nil
	case key.Matches(msg, layer.Quit), key.Matches(msg, m.Keymap.Cancel):
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.Keymap.Op = m.op
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case common.CloseViewMsg:
		m.Keymap.ResetMode()
		m.op = &operations.Noop{}
		return m, nil
	case confirmation.CloseMsg:
		if m.confirmation != nil {
			m.confirmation = nil
			return m, nil
		}
	case common.SetOperationMsg:
		m.op = msg.Operation
		return m, nil
	case common.UpdateRevSetMsg:
		m.revsetModel.Value = string(msg)
		if selectedRevision := m.selectedRevision(); selectedRevision != nil {
			cmds = append(cmds, common.Refresh(selectedRevision.GetChangeId()))
		} else {
			cmds = append(cmds, common.Refresh("@"))
		}
	case common.RefreshMsg:
		cmds = append(cmds,
			tea.Sequence(
				m.FetchRevisions(m.revsetModel.Value),
				common.SelectRevision(msg.SelectedRevision),
			))
	case common.SelectRevisionMsg:
		r := string(msg)
		idx := slices.IndexFunc(m.rows, func(row jj.GraphRow) bool {
			if r == "@" {
				return row.Commit.IsWorkingCopy
			}
			return row.Commit.GetChangeId() == r
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
	case tea.WindowSizeMsg:
		m.Width = msg.Width
	}

	if m.confirmation != nil {
		var cmd tea.Cmd
		var cm confirmation.Model
		if cm, cmd = m.confirmation.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.confirmation = &cm
		return m, tea.Batch(cmds...)
	}

	if op, ok := m.op.(operations.OperationWithOverlay); ok {
		var cmd tea.Cmd
		if m.op, cmd = op.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	if m.revsetModel.Editing {
		var cmd tea.Cmd
		if m.revsetModel, cmd = m.revsetModel.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.Keymap.Down):
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		default:
			if op, ok := m.op.(operations.HandleKey); ok {
				cmd = op.HandleKey(msg)
			} else {
				switch m.Keymap.Current {
				case ' ':
					m, cmd = m.handleBaseKeys(msg)
				}
			}
		}
	}
	if op, ok := m.op.(operations.TracksSelectedRevision); ok {
		op.SetSelectedRevision(m.selectedRevision())
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	topView := m.revsetModel.View()
	content := ""

	if m.status == common.Loading {
		content = "loading"
	}

	if m.status == common.Error {
		content = fmt.Sprintf("error: %v", m.error)
	}

	if m.status == common.Ready {
		if m.confirmation != nil {
			topView = lipgloss.JoinVertical(0, topView, m.confirmation.View())
		}
		h := m.Height - lipgloss.Height(topView)

		if m.viewRange.end-m.viewRange.start > h {
			m.viewRange.end = m.viewRange.start + h
		}

		var w jj.GraphWriter
		selectedLineStart := -1
		selectedLineEnd := -1
		for i, row := range m.rows {
			nodeRenderer := SegmentedRenderer{
				Palette:       common.DefaultPalette,
				op:            m.op,
				IsHighlighted: i == m.cursor,
			}

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

	return lipgloss.JoinVertical(0, topView, content)
}

func New(jj jj.Commands) Model {
	v := viewRange{start: 0, end: 0}
	defaultRevSet, _ := jj.GetConfig("revsets.log")
	return Model{
		status:      common.Loading,
		UICommands:  common.NewUICommands(jj),
		rows:        nil,
		viewRange:   &v,
		op:          &operations.Noop{},
		cursor:      0,
		Width:       20,
		revsetModel: revset.New(string(defaultRevSet)),
		Keymap:      keymap.NewKeyMap(),
	}
}
