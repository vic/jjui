package describe

import (
	"bytes"
	"github.com/idursun/jjui/internal/jj"
	"testing"
	"time"

	"github.com/idursun/jjui/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestCancel(t *testing.T) {
	c := test.NewTestContext(t)
	defer c.Verify()

	shell := test.NewShell(New(c, "revision", "description", 20))
	tm := teatest.NewTestModel(t, shell, teatest.WithInitialTermSize(100, 100))
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestEdit(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.Describe("revision", "description changed"))
	defer c.Verify()

	shell := test.NewShell(New(c, "revision", "description", 30))

	tm := teatest.NewTestModel(t, shell, teatest.WithInitialTermSize(100, 100))
	tm.Type(" changed")
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("changed"))
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
