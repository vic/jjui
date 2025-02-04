package jj

import (
	"bytes"
	"fmt"
	"os/exec"
)

type JJ struct {
	Location string
}

type Command interface {
	CombinedOutput() ([]byte, error)
	GetCommand() *exec.Cmd
	Args() []string
}

type commandWrapper struct {
	command *exec.Cmd
}

func (c commandWrapper) GetCommand() *exec.Cmd {
	return c.command
}

func (c commandWrapper) CombinedOutput() ([]byte, error) {
	return c.command.CombinedOutput()
}

func (c commandWrapper) Args() []string {
	return c.command.Args
}

func (jj JJ) createCommand(name string, args ...string) commandWrapper {
	cmd := exec.Command(name, args...)
	cmd.Dir = jj.Location
	return commandWrapper{
		command: cmd,
	}
}

func (jj JJ) GetConfig(key string) ([]byte, error) {
	cmd := exec.Command("jj", "config", "get", key)
	cmd.Dir = jj.Location
	output, err := cmd.CombinedOutput()
	output = bytes.Trim(output, "\n")
	return output, err
}

func (jj JJ) RebaseCommand(from string, to string) Command {
	return jj.createCommand("jj", "rebase", "-r", from, "-d", to)
}

func (jj JJ) RebaseBranchCommand(from string, to string) Command {
	return jj.createCommand("jj", "rebase", "-b", from, "-d", to)
}

func (jj JJ) Squash(from string, destination string) Command {
	return jj.createCommand("jj", "squash", "--from", from, "--into", destination)
}

func (jj JJ) SetDescription(rev string, description string) Command {
	return jj.createCommand("jj", "describe", "-r", rev, "-m", description)
}

func (jj JJ) ListBookmark(revision string) Command {
	return jj.createCommand("jj", "log", "-r", fmt.Sprintf("::%s- & bookmarks()", revision), "--template", "local_bookmarks.map(|x| x.name() ++ '\n')", "--no-graph", "--color", "never")
}

func (jj JJ) SetBookmark(revision string, name string) Command {
	return jj.createCommand("jj", "bookmark", "set", "-r", revision, name)
}

func (jj JJ) MoveBookmark(revision string, bookmark string) Command {
	return jj.createCommand("jj", "bookmark", "move", bookmark, "--to", revision)
}

func (jj JJ) DeleteBookmark(bookmark string) Command {
	return jj.createCommand("jj", "bookmark", "delete", bookmark)
}

func (jj JJ) GitFetch() Command {
	return jj.createCommand("jj", "git", "fetch")
}

func (jj JJ) GitPush() Command {
	return jj.createCommand("jj", "git", "push")
}

func (jj JJ) Diff(revision string, fileName string) Command {
	args := []string{"diff", "-r", revision, "--color", "always"}
	if fileName != "" {
		args = append(args, fileName)
	}

	return jj.createCommand("jj", args...)
}

func (jj JJ) Edit(revision string) Command {
	return jj.createCommand("jj", "edit", "-r", revision)
}

func (jj JJ) DiffEdit(revision string) Command {
	return jj.createCommand("jj", "diffedit", "-r", revision)
}

func (jj JJ) Split(revision string) Command {
	return jj.createCommand("jj", "split", "-r", revision)
}

func (jj JJ) Abandon(revision string) Command {
	return jj.createCommand("jj", "abandon", "-r", revision)
}

func (jj JJ) New(from string) Command {
	return jj.createCommand("jj", "new", "-r", from)
}

func (jj JJ) Status(revision string) Command {
	return jj.createCommand("jj", "log", "-r", revision, "--summary", "--no-graph", "--color", "never", "--template", "")
}

func (jj JJ) Restore(revision string, files []string) Command {
	args := []string{"restore", "-c", revision}
	args = append(args, files...)
	return jj.createCommand("jj", args...)
}
