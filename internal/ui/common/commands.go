package common

import (
	"github.com/charmbracelet/bubbletea"
	"jjui/internal/jj"
	"os/exec"
)

type Commands struct {
	jj jj.JJCommands
}

func (c Commands) GitFetch() tea.Cmd {
	f := func() tea.Msg {
		output, err := c.jj.GitFetch()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git fetch"), f)
}

func (c Commands) GitPush() tea.Cmd {
	f := func() tea.Msg {
		output, err := c.jj.GitPush()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git push"), f)
}

func (c Commands) Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := c.jj.RebaseCommand
	if operation == RebaseBranchOperation {
		rebase = c.jj.RebaseBranchCommand
	}
	output, err := rebase(from, to)
	return ShowOutput(string(output), err)
}

func (c Commands) SetDescription(revision string, description string) tea.Cmd {
	output, err := c.jj.SetDescription(revision, description)
	return ShowOutput(string(output), err)
}

func (c Commands) MoveBookmark(revision string, bookmark string) tea.Cmd {
	output, err := c.jj.MoveBookmark(revision, bookmark)
	return ShowOutput(string(output), err)
}

func (c Commands) FetchRevisions(revset string) tea.Cmd {
	return func() tea.Msg {
		graphLines, err := c.jj.GetCommits(revset)
		if err != nil {
			return UpdateRevisionsFailedMsg(err)
		}
		return UpdateRevisionsMsg(graphLines)
	}
}

func (c Commands) FetchBookmarks(revision string) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := c.jj.ListBookmark(revision)
		return UpdateBookmarksMsg(bookmarks)
	}
}

func (c Commands) SetBookmark(revision string, name string) tea.Cmd {
	output, err := c.jj.SetBookmark(revision, name)
	return ShowOutput(string(output), err)
}

func (c Commands) GetDiff(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := c.jj.Diff(revision)
		return ShowDiffMsg(output)
	}
}

func (c Commands) Edit(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := c.jj.Edit(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func (c Commands) DiffEdit(revision string) tea.Cmd {
	return tea.ExecProcess(exec.Command("jj", "diffedit", "-r", revision), func(err error) tea.Msg {
		return Refresh(revision)
	})
}

func (c Commands) Split(revision string) tea.Cmd {
	return tea.ExecProcess(exec.Command("jj", "split", "-r", revision), func(err error) tea.Msg {
		return Refresh(revision)
	})
}

func (c Commands) Abandon(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := c.jj.Abandon(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func (c Commands) NewRevision(from string) tea.Cmd {
	output, err := c.jj.New(from)
	return ShowOutput(string(output), err)
}

func NewCommands(jj jj.JJCommands) Commands {
	return Commands{jj}
}
