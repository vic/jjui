package oplog

import (
	"bytes"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
)

var normalStyle = lipgloss.NewStyle()

type updateOpLogMsg struct {
	Rows []Row
}

type viewRange struct {
	start int
	end   int
}
type Model struct {
	context   context.AppContext
	rows      []Row
	cursor    int
	keymap    config.KeyMappings[key.Binding]
	viewRange *viewRange
	width     int
	height    int
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

func (m *Model) Init() tea.Cmd {
	return m.load()
}

var restore = key.NewBinding(key.WithKeys("r"))

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateOpLogMsg:
		m.rows = msg.Rows
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Cancel):
			return m, common.Close
		case key.Matches(msg, m.keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			m.context.SetSelectedItem(context.SelectedOperation{OperationId: m.rows[m.cursor].OperationId})
			return m, common.SelectionChanged
		case key.Matches(msg, m.keymap.Down):
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
			m.context.SetSelectedItem(context.SelectedOperation{OperationId: m.rows[m.cursor].OperationId})
			return m, common.SelectionChanged
		case key.Matches(msg, m.keymap.OpLog.Restore):
			return m, tea.Batch(common.Close, m.context.RunCommand(jj.OpRestore(m.rows[m.cursor].OperationId), common.Refresh))
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.rows == nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "loading")
	}

	h := m.height
	viewHeight := m.viewRange.end - m.viewRange.start
	if viewHeight != h {
		m.viewRange.end = m.viewRange.start + h
	}
	var w Renderer
	w.Width = m.width
	selectedLineStart := -1
	selectedLineEnd := -1
	for i, row := range m.rows {
		isHighlighted := m.cursor == i
		if isHighlighted {
			selectedLineStart = w.LineCount()
		}
		w.RenderRow(row, isHighlighted)
		if isHighlighted {
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
	content = lipgloss.PlaceHorizontal(m.width, lipgloss.Left, content)
	return normalStyle.MaxWidth(m.width).Render(content)
}

func (m *Model) load() tea.Cmd {
	return func() tea.Msg {
		output, err := m.context.RunCommandImmediate(jj.OpLog())
		if err != nil {
			panic(err)
		}

		rows := ParseRows(bytes.NewReader(output))
		return updateOpLogMsg{Rows: rows}
	}
}

func New(context context.AppContext, width int, height int) *Model {
	keyMap := context.KeyMap()
	v := viewRange{start: 0, end: 0}
	return &Model{
		context:   context,
		keymap:    keyMap,
		rows:      nil,
		cursor:    0,
		viewRange: &v,
		width:     width,
		height:    height,
	}
}
