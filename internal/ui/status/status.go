package status

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/internal/ui/operations"
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

type clearMsg string

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
		if m.command == string(msg) {
			m.command = ""
			m.error = nil
			m.output = ""
		}
		return m, nil
	case common.CommandRunningMsg:
		m.command = string(msg)
		m.running = true
		return m, m.spinner.Tick
	case common.CommandCompletedMsg:
		m.running = false
		m.output = msg.Output
		m.error = msg.Err
		commandToBeCleared := m.command
		return m, tea.Tick(CommandClearDuration, func(time.Time) tea.Msg {
			return clearMsg(commandToBeCleared)
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
	s := common.DefaultPalette.StatusNormal.Render(" ")
	if m.running {
		s = common.DefaultPalette.StatusNormal.Render(m.spinner.View())
	} else if m.error != nil {
		s = common.DefaultPalette.StatusError.Render("✗ ")
	} else if m.command != "" {
		s = common.DefaultPalette.StatusSuccess.Render("✓ ")
	} else {
		if o, ok := m.op.(help.KeyMap); ok {
			s = m.help.View(o)
		}
	}
	ret := common.DefaultPalette.StatusNormal.Width(m.width - 2).SetString(m.command).Render()
	mode := common.DefaultPalette.StatusMode.Width(10).Render(m.op.Name())
	ret = lipgloss.JoinHorizontal(lipgloss.Left, mode, " ", s, ret)
	if m.error != nil {
		ret += " " + common.DefaultPalette.StatusError.Render(fmt.Sprintf("\n%v\n%s", m.error, m.output))
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
