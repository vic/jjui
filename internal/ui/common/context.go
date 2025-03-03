package common

import (
	"bytes"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/ui/operations"
	"os/exec"
)

type AppContext interface {
	Op() operations.Operation
	SetOp(op operations.Operation)
	KeyMap() KeyMappings[key.Binding]
	SelectedItem() SelectedItem
	SetSelectedItem(item SelectedItem)
	RunCommandImmediate(args []string) ([]byte, error)
	RunCommand(args []string, continuations ...tea.Cmd) tea.Cmd
	RunInteractiveCommand(args []string, continuation tea.Cmd) tea.Cmd
}

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
	config       Config
	op           operations.Operation
}

func (a *MainContext) Op() operations.Operation {
	return a.op
}

func (a *MainContext) SetOp(op operations.Operation) {
	a.op = op
}

func (a *MainContext) KeyMap() KeyMappings[key.Binding] {
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
			return CommandCompletedMsg{
				Output: string(output),
				Err:    err,
			}
		})
	commands = append(commands, continuations...)
	return tea.Batch(
		CommandRunning(args),
		tea.Sequence(commands...),
	)
}

func (a *MainContext) RunInteractiveCommand(args []string, continuation tea.Cmd) tea.Cmd {
	c := exec.Command("jj", args...)
	c.Dir = a.location
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return continuation()
	})
}

func NewAppContext(location string) AppContext {
	config := NewConfiguration()
	return &MainContext{
		location: location,
		config:   config,
	}
}
