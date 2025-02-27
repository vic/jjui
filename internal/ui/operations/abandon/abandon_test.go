package abandon

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/test"
	"testing"
	"time"
)

const revision = "revision"

func Test_Accept(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.Abandon(revision))
	defer c.Verify()

	model := test.NewShell(New(c, revision))
	tm := teatest.NewTestModel(t, model)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("abandon"))
	})

	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("closed"))
	})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func Test_Cancel(t *testing.T) {
	c := test.NewTestContext(t)
	defer c.Verify()

	model := test.NewShell(New(c, revision))
	tm := teatest.NewTestModel(t, model)
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("closed"))
	})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
