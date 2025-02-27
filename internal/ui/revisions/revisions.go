package revisions

import (
	"bytes"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
	"github.com/idursun/jjui/internal/ui/operations/abandon"
	"github.com/idursun/jjui/internal/ui/operations/bookmark"
	"github.com/idursun/jjui/internal/ui/operations/describe"
	"github.com/idursun/jjui/internal/ui/operations/details"
	"github.com/idursun/jjui/internal/ui/operations/git"
	"github.com/idursun/jjui/internal/ui/operations/rebase"
	"github.com/idursun/jjui/internal/ui/operations/squash"
	"github.com/idursun/jjui/internal/ui/operations/undo"
	"github.com/idursun/jjui/internal/ui/revset"
	"slices"
)

type viewRange struct {
	start int
	end   int
}

var normalStyle = lipgloss.NewStyle()

type Model struct {
	rows        []jj.GraphRow
	op          operations.Operation
	revsetValue string
	viewRange   *viewRange
	cursor      int
	width       int
	height      int
	context     common.AppContext
}

type updateRevisionsMsg struct {
	rows             []jj.GraphRow
	selectedRevision string
}

func (m *Model) IsFocused() bool {
	if _, ok := m.op.(common.Focusable); ok {
		return true
	}
	return false
}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return m.height
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetHeight(h int) {
	m.height = h
}

func (m *Model) ShortHelp() []key.Binding {
	if op, ok := m.op.(help.KeyMap); ok {
		return op.ShortHelp()
	}
	return (&operations.Noop{}).ShortHelp()
}

func (m *Model) FullHelp() [][]key.Binding {
	if op, ok := m.op.(help.KeyMap); ok {
		return op.FullHelp()
	}
	return [][]key.Binding{m.ShortHelp()}
}

func (m *Model) SelectedRevision() *jj.Commit {
	if m.cursor >= len(m.rows) {
		return nil
	}
	return m.rows[m.cursor].Commit
}

func (m *Model) Init() tea.Cmd {
	return common.Refresh
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	preSelectedRevision := m.SelectedRevision()
	switch msg := msg.(type) {
	case common.CloseViewMsg:
		m.op = &operations.Noop{}
		m.context.SetSelectedItem(common.SelectedRevision{ChangeId: m.SelectedRevision().ChangeId})
		return m, common.SelectionChanged
	case common.SetOperationMsg:
		m.op = msg.Operation
		return m, nil
	case revset.UpdateRevSetMsg:
		m.revsetValue = string(msg)
		if selectedRevision := m.SelectedRevision(); selectedRevision != nil {
			cmds = append(cmds, common.Refresh)
		} else {
			cmds = append(cmds, common.Refresh)
		}
	case common.RefreshMsg:
		return m, m.load(m.revsetValue, msg.SelectedRevision)
	case updateRevisionsMsg:
		if msg.rows != nil {
			currentSelectedRevision := msg.selectedRevision
			if cur := m.SelectedRevision(); currentSelectedRevision == "" && cur != nil {
				currentSelectedRevision = cur.GetChangeId()
			}
			m.rows = msg.rows
			m.cursor = m.selectRevision(currentSelectedRevision)
			if m.cursor == -1 {
				m.cursor = m.selectRevision("@")
			}
			if m.cursor == -1 {
				m.cursor = 0
			}
		}
	}

	if op, ok := m.op.(operations.OperationWithOverlay); ok {
		var cmd tea.Cmd
		if m.op, cmd = op.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, operations.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, operations.Down):
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		default:
			if op, ok := m.op.(operations.HandleKey); ok {
				cmd = op.HandleKey(msg)
			} else {
				switch {
				case key.Matches(msg, operations.Cancel):
					m.op = &operations.Noop{}
				case key.Matches(msg, operations.Details):
					m.op, cmd = details.NewOperation(m.context, m.SelectedRevision())
				case key.Matches(msg, operations.Undo):
					m.op, cmd = undo.NewOperation(m.context)
				case key.Matches(msg, operations.New):
					cmd = m.context.RunCommand(jj.New(m.SelectedRevision().GetChangeId()), common.RefreshAndSelect("@"))
				case key.Matches(msg, operations.Edit):
					cmd = m.context.RunCommand(jj.Edit(m.SelectedRevision().GetChangeId()), common.Refresh)
				case key.Matches(msg, operations.Diffedit):
					changeId := m.SelectedRevision().GetChangeId()
					cmd = m.context.RunInteractiveCommand(jj.DiffEdit(changeId), common.Refresh)
				case key.Matches(msg, operations.Abandon):
					m.op, cmd = abandon.NewOperation(m.context, m.SelectedRevision())
				case key.Matches(msg, operations.Split):
					currentRevision := m.SelectedRevision().GetChangeId()
					return m, m.context.RunInteractiveCommand(jj.Split(currentRevision, []string{}), common.Refresh)
				case key.Matches(msg, operations.Description):
					m.op, cmd = describe.NewOperation(m.context, m.SelectedRevision(), m.Width())
				case key.Matches(msg, operations.Diff):
					return m, func() tea.Msg {
						output, _ := m.context.RunCommandImmediate(jj.Diff(m.SelectedRevision().GetChangeId(), ""))
						return common.ShowDiffMsg(output)
					}
				case key.Matches(msg, operations.Refresh):
					cmd = common.Refresh
				case key.Matches(msg, operations.GitMode):
					m.op = git.NewOperation(m.context)
				case key.Matches(msg, operations.SquashMode):
					m.op = squash.NewOperation(m.context, m.SelectedRevision().ChangeIdShort)
					if m.cursor < len(m.rows)-1 {
						m.cursor++
					}
				case key.Matches(msg, operations.RebaseMode):
					m.op = rebase.NewOperation(m.context, m.SelectedRevision().ChangeIdShort, rebase.SourceRevision, rebase.TargetDestination)
				case key.Matches(msg, operations.BookmarkMode):
					m.op = bookmark.NewChooseBookmarkOperation(m.context)
				case key.Matches(msg, operations.Quit), key.Matches(msg, operations.Cancel):
					return m, tea.Quit
				}
			}
		}
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	if op, ok := m.op.(operations.TracksSelectedRevision); ok {
		op.SetSelectedRevision(m.SelectedRevision())
	}
	curSelected := m.SelectedRevision()
	if preSelectedRevision != curSelected {
		cmds = append(cmds, func() tea.Msg {
			return common.SelectionChangedMsg{}
		})
		m.context.SetSelectedItem(common.SelectedRevision{ChangeId: curSelected.ChangeId})
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if len(m.rows) == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "loading")
	}
	topView := ""
	topViewHeight := 0

	if m.op.RenderPosition() == operations.RenderPositionTop {
		topView = m.op.Render()
		topViewHeight = lipgloss.Height(topView)
	}

	h := m.height - topViewHeight
	viewHeight := m.viewRange.end - m.viewRange.start
	if viewHeight != h {
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

	content := w.String(m.viewRange.start, m.viewRange.end)
	content = lipgloss.PlaceHorizontal(m.Width(), lipgloss.Left, content)

	if topViewHeight > 0 {
		return lipgloss.JoinVertical(0, topView, content)
	}
	return normalStyle.MaxWidth(m.width).Render(content)
}

func (m *Model) load(revset string, selectedRevision string) tea.Cmd {
	return func() tea.Msg {
		output, err := m.context.RunCommandImmediate(jj.Log(revset))
		if err != nil {
			return common.UpdateRevisionsFailedMsg(err)
		}
		p := jj.NewParser(bytes.NewReader(output))
		graphLines := p.Parse()
		return updateRevisionsMsg{graphLines, selectedRevision}
	}
}

func (m *Model) selectRevision(revision string) int {
	idx := slices.IndexFunc(m.rows, func(row jj.GraphRow) bool {
		if revision == "@" {
			return row.Commit.IsWorkingCopy
		}
		return row.Commit.GetChangeId() == revision || row.Commit.ChangeIdShort == revision
	})
	return idx
}

func New(c common.AppContext) Model {
	v := viewRange{start: 0, end: 0}
	return Model{
		context:   c,
		rows:      nil,
		viewRange: &v,
		op:        &operations.Noop{},
		cursor:    0,
		width:     20,
		height:    10,
	}
}
