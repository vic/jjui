package common

import (
	"os"
	"os/exec"

	"jjui/internal/jj"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	CloseViewMsg        struct{}
	SelectRevisionMsg   string
	ShowDiffMsg         string
	UpdateRevisionsMsg  []jj.GraphRow
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
	RebaseRevision
	RebaseBranch
)

func Close() tea.Msg {
	return CloseViewMsg{}
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
	return tea.Sequence(CommandRunning("jj git fetch"), f, FetchRevisions(os.Getenv("PWD")))
}

func GitPush() tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.GitPush()
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(CommandRunning("jj git push"), f, FetchRevisions(os.Getenv("PWD")))
}

func Rebase(from, to string, operation Operation) tea.Cmd {
	rebase := jj.RebaseCommand
	if operation == RebaseBranch {
		rebase = jj.RebaseBranchCommand
	}
	output, err := rebase(from, to)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")), SelectRevision(from))
}

func SetDescription(revision string, description string) tea.Cmd {
	output, err := jj.SetDescription(revision, description)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
}

func MoveBookmark(revision string, bookmark string) tea.Cmd {
	output, err := jj.MoveBookmark(revision, bookmark)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")))
}

func FetchRevisions(location string) tea.Cmd {
	return func() tea.Msg {
		rows := jj.GetCommits(location)
		return UpdateRevisionsMsg(rows)
	}
}

func FetchBookmarks(revision string) tea.Cmd {
	return func() tea.Msg {
		bookmarks, _ := jj.ListBookmark(revision)
		return UpdateBookmarksMsg(bookmarks)
	}
}

func GetDiff(revision string) tea.Cmd {
	return func() tea.Msg {
		output, _ := jj.Diff(revision)
		return ShowDiffMsg(output)
	}
}

func Edit(revision string) tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.Edit(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")), SelectRevision("@"))
}

func DiffEdit(revision string) tea.Cmd {
	return tea.Sequence(tea.ExecProcess(exec.Command("jj", "diffedit", "-r", revision), nil), tea.ClearScreen)
}

func Split(revision string) tea.Cmd {
	return tea.Sequence(tea.ExecProcess(exec.Command("jj", "split", "-r", revision), nil), tea.ClearScreen, FetchRevisions(os.Getenv("PWD")))
}

func Abandon(revision string) tea.Cmd {
	f := func() tea.Msg {
		output, err := jj.Abandon(revision)
		return CommandCompletedMsg{Output: string(output), Err: err}
	}
	return tea.Sequence(f, FetchRevisions(os.Getenv("PWD")))
}

func NewRevision(from string) tea.Cmd {
	output, err := jj.New(from)
	return tea.Sequence(ShowOutput(string(output), err), FetchRevisions(os.Getenv("PWD")), SelectRevision("@"))
}
