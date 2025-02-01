package details

import (
	"bytes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"jjui/internal/ui/common"
	"jjui/test"
	"testing"
	"time"
)

const (
	Revision     = "ignored"
	StatusOutput = "M file.txt\nA newfile.txt\n"
)

func TestModel_Init_ExecutesStatusCommand(t *testing.T) {
	commands := test.NewJJCommands()
	statusCommand := commands.ExpectStatus(t, Revision)
	statusCommand.Output = []byte(StatusOutput)

	tm := teatest.NewTestModel(t, New(Revision, common.NewUICommands(commands)))
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
	commands.Verify(t)
}

func TestModel_Update_RestoresSelectedFiles(t *testing.T) {
	commands := test.NewJJCommands()
	statusCommand := commands.ExpectStatus(t, Revision)
	statusCommand.Output = []byte(StatusOutput)

	commands.ExpectRestore(t, Revision, []string{"file.txt"})

	tm := teatest.NewTestModel(t, New(Revision, common.NewUICommands(commands)))
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})

	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Yes"))
	})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
	commands.Verify(t)
}
