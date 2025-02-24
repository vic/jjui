package undo

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/confirmation"
	"github.com/idursun/jjui/internal/ui/operations"
	"github.com/idursun/jjui/test"
	"testing"
	"time"
)

type OperationHost struct {
	closed    bool
	Operation operations.Operation
}

func (o OperationHost) Init() tea.Cmd {
	return nil
}

func (o OperationHost) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.CloseViewMsg:
		o.closed = true
		return o, nil
	case confirmation.CloseMsg:
		o.closed = true
		return o, nil
	default:
		if op, ok := o.Operation.(operations.OperationWithOverlay); ok {
			var cmd tea.Cmd
			o.Operation, cmd = op.Update(msg)
			return o, cmd
		}
	}
	return o, nil
}

func (o OperationHost) View() string {
	if o.closed {
		return "closed"
	}
	return o.Operation.Render()
}

func TestConfirm(t *testing.T) {
	commands := test.NewJJCommands(t)
	commands.ExpectUndo()
	defer commands.Verify()

	operation, _ := NewOperation(common.NewUICommands(commands))
	model := OperationHost{Operation: operation}

	tm := teatest.NewTestModel(t, model)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("undo"))
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("closed"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestCancel(t *testing.T) {
	commands := test.NewJJCommands(t)
	defer commands.Verify()

	operation, _ := NewOperation(common.NewUICommands(commands))
	model := OperationHost{Operation: operation}

	tm := teatest.NewTestModel(t, model)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("undo"))
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("closed"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
