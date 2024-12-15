package common

import (
	"os/exec"

	"jjui/internal/jj"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	CloseViewMsg        struct{}
	RefreshMsg          struct{}
	SelectRevisionMsg   string
	ShowDiffMsg         string
	UpdateRevisionsMsg  *jj.Dag
	UpdateBookmarksMsg  []jj.Bookmark
	CommandRunningMsg   string
	CommandCompletedMsg struct {
		Output string
		Err    error
	}
)

type Operation int

const (
	None Operation = iota
	RebaseRevisionOperation
	RebaseBranchOperation
	EditDescriptionOperation
	SetBookmarkOperation
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func Refresh() tea.Msg {
	return RefreshMsg{}
}

func SelectRevision(revision string) tea.Cmd {
	return func() tea.Msg {
		return SelectRevisionMsg(revision)
	}
}

func CommandRunning(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandRunningMsg(command)
	}
}

func ShowOutput(output string, err error) tea.Cmd {
	return func() tea.Msg {
		return CommandCompletedMsg{
			Output: output,
			Err:    err,
		}
	}
}

func GitFetch() tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.GitFetch()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git fetch"), f)
}

func GitPush() tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.GitPush()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git push"), f)
}

func Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := jj.RebaseCommand
	if operation == RebaseBranchOperation {
		rebase = jj.RebaseBranchCommand
	}
	output, err := rebase(from, to)
	return ShowOutput(string(output), err)
}

func SetDescription(revision string, description string) tea.Cmd {
	output, err := jj.SetDescription(revision, description)
	return ShowOutput(string(output), err)
}

func MoveBookmark(revision string, bookmark string) tea.Cmd {
	output, err := jj.MoveBookmark(revision, bookmark)
	return ShowOutput(string(output), err)
}

func FetchRevisions(location string, revset string) tea.Cmd {
	return func() tea.Msg {
		dag := jj.GetCommits(location, revset)
		return UpdateRevisionsMsg(dag)
	}
}

func FetchBookmarks(revision string) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := jj.ListBookmark(revision)
		return UpdateBookmarksMsg(bookmarks)
	}
}

func SetBookmark(revision string, name string) tea.Cmd {
	output, err := jj.SetBookmark(revision, name)
	return ShowOutput(string(output), err)
}

func GetDiff(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := jj.Diff(revision)
		return ShowDiffMsg(output)
	}
}

func Edit(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := jj.Edit(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func DiffEdit(revision string) tea.Cmd {
	return tea.Sequence(tea.ExecProcess(exec.Command("jj", "diffedit", "-r", revision), nil), tea.ClearScreen)
}

func Split(revision string) tea.Cmd {
	return tea.Sequence(tea.ExecProcess(exec.Command("jj", "split", "-r", revision), nil), tea.ClearScreen)
}

func Abandon(revision string) tea.Cmd {
	return func() tea.Msg {
		output, err := jj.Abandon(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
}

func NewRevision(from string) tea.Cmd {
	output, err := jj.New(from)
	return ShowOutput(string(output), err)
}
