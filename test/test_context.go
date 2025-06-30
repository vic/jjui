package test

import (
	"bytes"
	"context"
	appContext "github.com/idursun/jjui/internal/ui/context"
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/stretchr/testify/assert"
)

type ExpectedCommand struct {
	args   []string
	output []byte
	called bool
}

func (e *ExpectedCommand) SetOutput(output []byte) *ExpectedCommand {
	e.output = output
	return e
}

type TestCommandRunner struct {
	*testing.T
	expectations map[string][]*ExpectedCommand
}

func (t *TestCommandRunner) RunCommandImmediate(args []string) ([]byte, error) {
	subCommand := args[0]
	if _, ok := t.expectations[subCommand]; !ok {
		assert.Fail(t, "unexpected command", subCommand)
	}
	expectations := t.expectations[subCommand]
	if len(expectations) == 0 {
		assert.Fail(t, "unexpected command", subCommand)
	}
	for _, e := range expectations {
		if assert.Equal(t.T, e.args, args) {
			e.called = true
			return e.output, nil
		}
	}
	assert.Fail(t, "unexpected command", subCommand)
	return nil, nil
}

func (t *TestCommandRunner) RunCommandStreaming(_ context.Context, args []string) (*appContext.StreamingCommand, error) {
	reader, err := t.RunCommandImmediate(args)
	return &appContext.StreamingCommand{
		ReadCloser: io.NopCloser(bytes.NewReader(reader)),
		ErrPipe:    nil,
	}, err
}

func (t *TestCommandRunner) RunCommand(args []string, continuations ...tea.Cmd) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	cmds = append(cmds, func() tea.Msg {
		_, _ = t.RunCommandImmediate(args)
		return common.CommandCompletedMsg{}
	})
	cmds = append(cmds, continuations...)
	return tea.Batch(cmds...)
}

func (t *TestCommandRunner) RunInteractiveCommand(args []string, continuation tea.Cmd) tea.Cmd {
	return t.RunCommand(args, continuation)
}

func (t *TestCommandRunner) Expect(args []string) *ExpectedCommand {
	subCommand := args[0]
	if _, ok := t.expectations[subCommand]; !ok {
		t.expectations[subCommand] = make([]*ExpectedCommand, 0)
	}
	e := &ExpectedCommand{
		args: args,
	}
	t.expectations[subCommand] = append(t.expectations[subCommand], e)
	return e
}

func (t *TestCommandRunner) Verify() {
	for subCommand, expectations := range t.expectations {
		for _, e := range expectations {
			if !e.called {
				assert.Fail(t, "expected command not called", subCommand)
			}
		}
	}
}

func NewTestCommandRunner(t *testing.T) *TestCommandRunner {
	return &TestCommandRunner{
		T:            t,
		expectations: make(map[string][]*ExpectedCommand),
	}
}

func NewTestContext(commandRunner appContext.CommandRunner) *appContext.MainContext {
	return &appContext.MainContext{
		CommandRunner: commandRunner,
		SelectedItem:  nil,
	}
}
