package details

import (
	"bytes"
	"testing"
	"time"

	"github.com/idursun/jjui/test"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/idursun/jjui/internal/ui/common"
)

const (
	Revision     = "ignored"
	StatusOutput = "M file.txt\nA newfile.txt\n"
)

func TestModel_Init_ExecutesStatusCommand(t *testing.T) {
	commands := test.NewJJCommands(t)
	defer commands.Verify()

	statusCommand := commands.ExpectStatus(t, Revision)
	statusCommand.Output = []byte(StatusOutput)

	tm := teatest.NewTestModel(t, New(Revision, common.NewUICommands(commands)))
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})
	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestModel_Update_RestoresSelectedFiles(t *testing.T) {
	commands := test.NewJJCommands(t)
	defer commands.Verify()

	commands.ExpectStatus(t, Revision).Output = []byte(StatusOutput)

	commands.ExpectRestore(Revision, []string{"file.txt"})

	tm := teatest.NewTestModel(t, test.NewShell(New(Revision, common.NewUICommands(commands))))
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})

	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

func TestModel_Update_SplitsSelectedFiles(t *testing.T) {
	commands := test.NewJJCommands(t)
	defer commands.Verify()

	statusCommand := commands.ExpectStatus(t, Revision)
	statusCommand.Output = []byte(StatusOutput)

	commands.ExpectSplit(Revision, []string{"file.txt"})

	tm := teatest.NewTestModel(t, test.NewShell(New(Revision, common.NewUICommands(commands))))
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})

	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
