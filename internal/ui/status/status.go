package status

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
	"time"
)

type Model struct {
	context context.AppContext
	spinner spinner.Model
	help    help.Model
	command string
	running bool
	output  string
	error   error
	width   int
	op      operations.Operation
}

const CommandClearDuration = 3 * time.Second

var (
	normalStyle  = lipgloss.NewStyle()
	successStyle = lipgloss.NewStyle().Inherit(normalStyle).Foreground(common.Green)
	errorStyle   = lipgloss.NewStyle().Inherit(normalStyle).Foreground(common.Red)
	modeStyle    = lipgloss.NewStyle().Inherit(normalStyle).Foreground(common.Black).Background(common.DarkBlue)
)

type clearMsg struct{}

func (m *Model) Width() int {
	return m.width
}

func (m *Model) Height() int {
	return 1
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

func (m *Model) SetHeight(int) {}
func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clearMsg:
		m.command = ""
		m.error = nil
		m.output = ""
		return m, nil
	case common.CommandRunningMsg:
		m.command = string(msg)
		m.running = true
		return m, m.spinner.Tick
	case common.CommandCompletedMsg:
		m.running = false
		m.output = msg.Output
		m.error = msg.Err
		return m, tea.Tick(CommandClearDuration, func(time.Time) tea.Msg {
			return clearMsg{}
		})
	case operations.OperationChangedMsg:
		m.op = msg.Operation
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *Model) View() string {
	s := normalStyle.Render(" ")
	if m.running {
		s = normalStyle.Render(m.spinner.View())
	} else if m.error != nil {
		s = errorStyle.Render("✗ ")
	} else if m.command != "" {
		s = successStyle.Render("✓ ")
	} else {
		if o, ok := m.op.(help.KeyMap); ok {
			s = m.help.View(o)
		}
	}
	ret := normalStyle.Width(m.width - 2).SetString(m.command).Render()
	mode := modeStyle.Width(10).Render(m.op.Name())
	ret = lipgloss.JoinHorizontal(lipgloss.Left, mode, " ", s, ret)
	if m.error != nil {
		ret += " " + errorStyle.Render(fmt.Sprintf("\n%v\n%s", m.error, m.output))
	}
	return ret
}

func New(context context.AppContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	h := help.New()
	h.Styles.ShortKey = common.DefaultPalette.CommitShortStyle
	h.Styles.ShortDesc = common.DefaultPalette.CommitIdRestStyle
	h.ShortSeparator = " "
	return Model{
		context: context,
		op:      operations.Default(context),
		spinner: s,
		help:    h,
		command: "",
		running: false,
		output:  "",
	}
}
