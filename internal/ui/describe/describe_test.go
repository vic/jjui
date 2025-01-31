package describe

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"jjui/internal/jj"
	"jjui/internal/ui/common"
	"os/exec"
	"testing"
)

type TestJJCommands struct {
	Invocations [][]string
}

func (t *TestJJCommands) GetConfig(key string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) RebaseCommand(from string, to string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) RebaseBranchCommand(from string, to string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) SetDescription(rev string, description string) jj.Command {
	t.Invocations = append(t.Invocations, []string{"SetDescription", rev, description})
	return &exec.Cmd{}
}

func (t *TestJJCommands) ListBookmark(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) SetBookmark(revision string, name string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) MoveBookmark(revision string, bookmark string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) DeleteBookmark(bookmark string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) GitFetch() jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) GitPush() jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Diff(revision string, fineName string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Edit(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) DiffEdit(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Abandon(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) New(from string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Split(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) GetCommits(revset string) ([]jj.GraphRow, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Squash(from string, destination string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Status(revision string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func (t *TestJJCommands) Restore(revision string, files []string) jj.Command {
	//TODO implement me
	panic("implement me")
}

func TestCancel(t *testing.T) {
	model := New(common.NewCommands(&TestJJCommands{}), "revision", "description", 20)
	var cmd tea.Cmd
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.NotNil(t, cmd)
	msg := cmd()
	assert.Equal(t, common.CloseViewMsg{}, msg)
	assert.Equal(t, "revision", model.(Model).revision)
}

func TestEdit(t *testing.T) {
	commands := TestJJCommands{}
	model := New(common.NewCommands(&commands), "revision", "description", 20)
	var cmd tea.Cmd
	model, cmd = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, cmd)
	msg := cmd()
	assert.NotNil(t, msg)
	//assert.Equal(t, [][]string{{"SetDescription", "revision", "description"}}, commands.Invocations)
	assert.Contains(t, commands.Invocations, []string{"SetDescriptions", "revision", "description"})
}
