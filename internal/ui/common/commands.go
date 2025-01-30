package common

import (
	"github.com/charmbracelet/bubbletea"
	"jjui/internal/jj"
	"strings"
)

type Commands struct {
	jj jj.Commands
}

func (c Commands) GitFetch(revision string) tea.Cmd {
	return ShowOutput(c.jj.GitFetch(), Refresh(revision))
}

func (c Commands) GitPush(revision string) tea.Cmd {
	return ShowOutput(c.jj.GitPush(), Refresh(revision))
}

func (c Commands) Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := c.jj.RebaseCommand
	if operation == RebaseBranchOperation {
		rebase = c.jj.RebaseBranchCommand
	}
	cmd := rebase(from, to)
	return ShowOutput(cmd, Refresh(to))
}

func (c Commands) Squash(from, destination string) tea.Cmd {
	cmd := c.jj.Squash(from, destination)
	return tea.Sequence(
		CommandRunning(strings.Join(cmd.Args, " ")),
		tea.ExecProcess(cmd, func(err error) tea.Msg {
			return CommandCompletedMsg{Output: "", Err: err}
		}),
		Refresh(destination),
	)
}

func (c Commands) SetDescription(revision string, description string) tea.Cmd {
	return ShowOutput(c.jj.SetDescription(revision, description), Refresh(revision), Close)
}

func (c Commands) MoveBookmark(revision string, bookmark string) tea.Cmd {
	return ShowOutput(c.jj.MoveBookmark(revision, bookmark), Refresh(revision), Close)
}

func (c Commands) DeleteBookmark(revision, bookmark string) tea.Cmd {
	return ShowOutput(c.jj.DeleteBookmark(bookmark), Refresh(revision), Close)
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
		cmd := c.jj.ListBookmark(revision)
		//TODO: handle error
		output, _ := cmd.CombinedOutput()
		var bookmarks []string
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			bookmarks = append(bookmarks, line)
		}
		return UpdateBookmarksMsg{
			Bookmarks: bookmarks,
			Revision:  revision,
		}
	}
}

func (c Commands) SetBookmark(revision string, name string) tea.Cmd {
	return ShowOutput(c.jj.SetBookmark(revision, name), Refresh(revision), Close)
}

func (c Commands) GetDiff(revision string, fileName string) tea.Cmd {
	return func() tea.Msg {
		output, _ := c.jj.Diff(revision, fileName).CombinedOutput()
		return ShowDiffMsg(output)
	}
}

func (c Commands) Edit(revision string) tea.Cmd {
	return ShowOutput(c.jj.Edit(revision), Refresh("@"))
}

func (c Commands) DiffEdit(revision string) tea.Cmd {
	return tea.ExecProcess(c.jj.DiffEdit(revision), func(err error) tea.Msg {
		return RefreshMsg{SelectedRevision: revision}
	})
}

func (c Commands) Split(revision string) tea.Cmd {
	return tea.ExecProcess(c.jj.Split(revision), func(err error) tea.Msg {
		return RefreshMsg{SelectedRevision: revision}
	})
}

func (c Commands) Abandon(revision string) tea.Cmd {
	return ShowOutput(c.jj.Abandon(revision), Refresh("@"), Close)
}

func (c Commands) NewRevision(from string) tea.Cmd {
	return ShowOutput(c.jj.New(from), Refresh("@"))
}

func (c Commands) Status(revision string) tea.Cmd {
	cmd := c.jj.Status(revision)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return func() tea.Msg {
			summary := strings.Split(strings.TrimSpace(string(output)), "\n")
			return UpdateCommitStatusMsg(summary)
		}
	}
	return func() tea.Msg {
		return CommandCompletedMsg{
			Output: string(output),
			Err:    err,
		}
	}
}

func NewCommands(jj jj.Commands) Commands {
	return Commands{jj}
}
