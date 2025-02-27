package bookmark

import (
	"bytes"
	"github.com/idursun/jjui/internal/jj"
	"testing"
	"time"

	"github.com/idursun/jjui/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestSetBookmarkModel_Update(t *testing.T) {
	c := test.NewTestContext(t)
	defer c.Verify()

	c.Expect(jj.BookmarkSet("revision", "name"))
	tm := teatest.NewTestModel(t, NewSetBookmark(c, "revision"))
	tm.Type("name")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("name"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
