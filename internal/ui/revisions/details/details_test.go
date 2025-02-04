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

type model struct {
	details tea.Model
}

func (m model) Init() tea.Cmd {
	return m.details.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(common.CloseViewMsg); ok {
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.details, cmd = m.details.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.details.View()
}

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

func TestModel_Update_SplitsSelectedFiles(t *testing.T) {
	commands := test.NewJJCommands()
	statusCommand := commands.ExpectStatus(t, Revision)
	statusCommand.Output = []byte(StatusOutput)

	commands.ExpectSplit(t, Revision, []string{"file.txt"})
	shell := model{
		details: New(Revision, common.NewUICommands(commands)),
	}

	tm := teatest.NewTestModel(t, shell)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("file.txt"))
	})

	tm.Send(tea.KeyMsg{Type: tea.KeySpace})
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
	commands.Verify(t)
}
