package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/idursun/jjui/internal/ui/operations"
	"github.com/idursun/jjui/internal/ui/preview"
	"github.com/idursun/jjui/internal/ui/revset"
	"strings"

	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/diff"
	"github.com/idursun/jjui/internal/ui/revisions"
	"github.com/idursun/jjui/internal/ui/status"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	TogglePreview = key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "toggle preview"))
	ToggleHelp    = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
)

type Model struct {
	revisions      tea.Model
	revsetModel    revset.Model
	previewModel   tea.Model
	helpVisible    bool
	previewVisible bool
	diff           tea.Model
	help           help.Model
	state          common.State
	error          error
	status         tea.Model
	output         string
	width          int
	height         int
	context        *common.AppContext
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(tea.SetWindowTitle("jjui"), m.revisions.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok && m.diff != nil {
		m.diff = nil
		return m, nil
	}

	var cmd tea.Cmd
	if m.diff != nil {
		m.diff, cmd = m.diff.Update(msg)
		return m, cmd
	}

	if m.revsetModel.Editing {
		var cmd tea.Cmd
		if m.revsetModel, cmd = m.revsetModel.Update(msg); cmd != nil {
			return m, cmd
		}
	}

	if r, ok := m.revisions.(common.Editable); ok && r.IsEditing() {
		m.revisions, cmd = m.revisions.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, operations.Cancel) && m.state == common.Error:
			m.state = common.Ready
			m.error = nil
		case key.Matches(msg, operations.Revset):
			m.revsetModel, _ = m.revsetModel.Update(revset.EditRevSetMsg{})
		case key.Matches(msg, ToggleHelp):
			return m, common.ToggleHelp
		case key.Matches(msg, TogglePreview):
			m.previewVisible = !m.previewVisible
		}
	case common.ToggleHelpMsg:
		m.helpVisible = !m.helpVisible
	case common.ShowDiffMsg:
		m.diff = diff.New(string(msg), m.width, m.height)
		return m, m.diff.Init()
	case common.CommandCompletedMsg:
		m.output = msg.Output
	case common.UpdateRevisionsFailedMsg:
		m.state = common.Error
		m.error = msg
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if r, ok := m.revisions.(common.Sizable); ok {
			if m.previewVisible {
				r.SetWidth(m.width / 2)
			} else {
				r.SetWidth(m.width)
			}
			r.SetHeight(m.height - 4)
		}
		if p, ok := m.previewModel.(common.Sizable); ok && m.previewVisible {
			p.SetWidth(m.width / 2)
			p.SetHeight(m.height - 4)
		}
		if s, ok := m.status.(common.Sizable); ok {
			s.SetWidth(m.width)
		}
	}

	m.revisions, cmd = m.revisions.Update(msg)

	var statusCmd tea.Cmd
	m.status, statusCmd = m.status.Update(msg)

	var previewCmd tea.Cmd
	m.previewModel, previewCmd = m.previewModel.Update(msg)
	return m, tea.Batch(cmd, statusCmd, previewCmd)
}

func (m Model) View() string {
	if m.diff != nil {
		return m.diff.View()
	}

	topView := m.revsetModel.View()

	if m.state == common.Error {
		topView += fmt.Sprintf("\nerror: %v\n", m.error)
	}

	topViewHeight := lipgloss.Height(topView)

	var b strings.Builder
	if h, ok := m.revisions.(help.KeyMap); ok && m.helpVisible {
		b.WriteString(m.help.View(h))
		b.WriteString("\n")
	}
	b.WriteString(m.status.View())

	footer := b.String()
	footerHeight := lipgloss.Height(footer)

	if r, ok := m.revisions.(common.Sizable); ok {
		r.SetWidth(m.width)
		if m.previewVisible {
			r.SetWidth(m.width / 2)
		}
		r.SetHeight(m.height - footerHeight - topViewHeight)
	}
	revisionsView := m.revisions.View()

	previewView := ""
	if p, ok := m.previewModel.(common.Sizable); ok && m.previewVisible {
		p.SetWidth(m.width - lipgloss.Width(revisionsView))
		p.SetHeight(m.height - footerHeight - topViewHeight)
		previewView = m.previewModel.View()
	}

	centerView := lipgloss.JoinHorizontal(lipgloss.Left, revisionsView, previewView)
	return lipgloss.JoinVertical(0, topView, centerView, footer)
}

func New(c *common.AppContext) tea.Model {
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	defaultRevSet, _ := c.JJ.GetConfig("revsets.log")
	revisionsModel := revisions.New(c)
	previewModel := preview.New(c)
	statusModel := status.New()
	return Model{
		context:      c,
		state:        common.Loading,
		revisions:    &revisionsModel,
		previewModel: &previewModel,
		help:         h,
		helpVisible:  true,
		status:       &statusModel,
		revsetModel:  revset.New(string(defaultRevSet)),
	}
}
