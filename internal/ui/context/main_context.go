package context

import (
	"bytes"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/ui/common"
	"os/exec"
)

type SelectedItem interface{}

type SelectedRevision struct {
	ChangeId string
}

type SelectedFile struct {
	ChangeId string
	File     string
}

type MainContext struct {
	selectedItem SelectedItem
	location     string
	config       *config.Config
}

func (a *MainContext) KeyMap() config.KeyMappings[key.Binding] {
	return a.config.GetKeyMap()
}

func (a *MainContext) SelectedItem() SelectedItem {
	return a.selectedItem
}

func (a *MainContext) SetSelectedItem(item SelectedItem) {
	a.selectedItem = item
}

func (a *MainContext) RunCommandImmediate(args []string) ([]byte, error) {
	c := exec.Command("jj", args...)
	c.Dir = a.location
	output, err := c.CombinedOutput()
	return bytes.Trim(output, "\n"), err
}

func (a *MainContext) RunCommand(args []string, continuations ...tea.Cmd) tea.Cmd {
	commands := make([]tea.Cmd, 0)
	commands = append(commands,
		func() tea.Msg {
			c := exec.Command("jj", args...)
			c.Dir = a.location
			output, err := c.CombinedOutput()
			return common.CommandCompletedMsg{
				Output: string(output),
				Err:    err,
			}
		})
	commands = append(commands, continuations...)
	return tea.Batch(
		common.CommandRunning(args),
		tea.Sequence(commands...),
	)
}

func (a *MainContext) RunInteractiveCommand(args []string, continuation tea.Cmd) tea.Cmd {
	c := exec.Command("jj", args...)
	errBuffer := &bytes.Buffer{}
	c.Stderr = errBuffer
	c.Dir = a.location
	return tea.Batch(
		common.CommandRunning(args),
		tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				return common.CommandCompletedMsg{Err: err, Output: errBuffer.String()}
			}
			return tea.Batch(continuation, func() tea.Msg {
				return common.CommandCompletedMsg{Err: nil}
			})()
		}),
	)
}

func NewAppContext(location string) AppContext {
	configuration := config.Load()
	return &MainContext{
		location: location,
		config:   configuration,
	}
}
