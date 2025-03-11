package git

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/test"
	"testing"
	"time"
)

func Test_Push(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.GitPush())
	defer c.Verify()

	op := NewModel(c, nil, 0, 0)
	tm := teatest.NewTestModel(t, test.NewShell(op))
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func Test_Fetch(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.GitFetch())
	defer c.Verify()

	op := NewModel(c, nil, 0, 0)
	tm := teatest.NewTestModel(t, test.NewShell(op))
	tm.Type("/")
	tm.Type("fetch")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
