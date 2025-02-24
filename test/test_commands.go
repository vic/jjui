package test

import (
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"

	"github.com/idursun/jjui/internal/jj"
)

type JJCommands struct {
	*testing.T
	expectations map[string]*MockedCommand
}

type MockedCommand struct {
	args   []string
	Output []byte
	Err    error
	called bool
}

func (m *MockedCommand) CombinedOutput() ([]byte, error) {
	m.called = true
	return m.Output, m.Err
}

func (m *MockedCommand) GetCommand() *exec.Cmd {
	m.called = true
	// do nothing
	return exec.Command(":")
}

func (m *MockedCommand) Args() []string {
	return m.args
}

type unexpectedCommand struct {
	name string
	t    *testing.T
}

func (u *unexpectedCommand) CombinedOutput() ([]byte, error) {
	assert.Failf(u.t, "unexpected call", u.name)
	return nil, nil
}

func (u *unexpectedCommand) GetCommand() *exec.Cmd {
	assert.Failf(u.t, "unexpected call", u.name)
	return nil
}

func (u *unexpectedCommand) Args() []string {
	return []string{}
}

func (t *JJCommands) GetConfig(key string) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) RebaseCommand(from string, to string, source string, target string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) SetDescription(rev string, description string) jj.Command {
	expectation := t.expectations["SetDescription"]
	if expectation == nil {
		return &unexpectedCommand{name: "SetDescription", t: t.T}
	}
	assert.Equal(t, expectation.args, []string{rev, description})
	return expectation
}

func (t *JJCommands) ListBookmark(revision string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) SetBookmark(revision string, name string) jj.Command {
	expectation := t.expectations["SetBookmark"]
	if expectation == nil {
		return &unexpectedCommand{name: "SetBookmark", t: t.T}
	}
	assert.Equal(t, expectation.args, []string{revision, name})
	return expectation
}

func (t *JJCommands) MoveBookmark(revision string, bookmark string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) DeleteBookmark(bookmark string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) GitFetch() jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) GitPush() jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Diff(revision string, fineName string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Edit(revision string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) DiffEdit(revision string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Abandon(revision string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) New(from string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Undo() jj.Command {
	expectation := t.expectations["Undo"]
	if expectation == nil {
		return &unexpectedCommand{name: "Undo", t: t.T}
	}
	return expectation
}

func (t *JJCommands) ExpectUndo() *MockedCommand {
	t.expectations["Undo"] = &MockedCommand{}
	return &MockedCommand{}
}

func (t *JJCommands) Split(revision string, files []string) jj.Command {
	expectation := t.expectations["Split"]
	if expectation == nil {
		return &unexpectedCommand{name: "Split", t: t.T}
	}
	assert.Equal(t, expectation.args, append([]string{revision}, files...))
	return expectation
}

func (t *JJCommands) ExpectSplit(revision string, files []string) *MockedCommand {
	command := MockedCommand{
		args: append([]string{revision}, files...),
	}
	t.expectations["Split"] = &command
	return &command
}

func (t *JJCommands) GetCommits(revset string) ([]jj.GraphRow, error) {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Squash(from string, destination string) jj.Command {
	// TODO implement me
	panic("implement me")
}

func (t *JJCommands) Status(revision string) jj.Command {
	expectation := t.expectations["Status"]
	if expectation == nil {
		return &unexpectedCommand{name: "Status", t: t.T}
	}
	assert.Equal(t, expectation.args, []string{revision})
	return expectation
}

func (t *JJCommands) ExpectStatus(tt *testing.T, revision string) *MockedCommand {
	command := MockedCommand{
		args: []string{revision},
	}
	t.expectations["Status"] = &command
	return &command
}

func (t *JJCommands) Restore(revision string, files []string) jj.Command {
	expectation := t.expectations["Restore"]
	if expectation == nil {
		return &unexpectedCommand{name: "Restore", t: t.T}
	}
	assert.Equal(t, expectation.args, append([]string{revision}, files...))
	return expectation
}

func (t *JJCommands) ExpectRestore(revision string, files []string) *MockedCommand {
	command := MockedCommand{
		args: append([]string{revision}, files...),
	}
	t.expectations["Restore"] = &command
	return &command
}

func (t *JJCommands) ExpectSetDescription(rev string, description string) *MockedCommand {
	command := MockedCommand{
		args: []string{rev, description},
	}
	t.expectations["SetDescription"] = &command
	return &command
}

func (t *JJCommands) ExpectSetBookmark(revision string, name string) *MockedCommand {
	command := MockedCommand{
		args: []string{revision, name},
	}
	t.expectations["SetBookmark"] = &command
	return &command
}

func (t *JJCommands) Verify() {
	for name, expectation := range t.expectations {
		if !expectation.called {
			t.Errorf("expected %s to be called", name)
		}
	}
}

func NewJJCommands(t *testing.T) *JJCommands {
	return &JJCommands{
		T:            t,
		expectations: make(map[string]*MockedCommand),
	}
}
