package git

import (
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/context"
	"github.com/idursun/jjui/test"
	"testing"
	"time"
)

func Test_Push(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.GitPush())
	c.SetSelectedItem(context.SelectedRevision{ChangeId: "revision"})
	defer c.Verify()

	op := NewOperation(c)
	tm := teatest.NewTestModel(t, test.OperationHost{Operation: op})
	tm.Type("p")
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func Test_Fetch(t *testing.T) {
	c := test.NewTestContext(t)
	c.Expect(jj.GitFetch())
	c.SetSelectedItem(context.SelectedRevision{ChangeId: "revision"})
	defer c.Verify()

	op := NewOperation(c)
	tm := teatest.NewTestModel(t, test.OperationHost{Operation: op})
	tm.Type("f")
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
