package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"jjui/internal/jj"
	"os/exec"
	"strings"
)

type (
	CloseViewMsg             struct{}
	RefreshMsg               struct{ SelectedRevision string }
	SelectRevisionMsg        string
	ShowDiffMsg              string
	UpdateRevSetMsg          string
	UpdateRevisionsMsg       []jj.GraphRow
	UpdateRevisionsFailedMsg error
	UpdateBookmarksMsg       struct {
		Bookmarks []string
		Revision  string
		Operation Operation
	}
	UpdateCommitStatusMsg []string
	CommandRunningMsg     string
	AbandonMsg            string
	SetDescriptionMsg     struct {
		Revision    string
		Description string
	}
	SetBookmarkMsg struct {
		Revision string
		Bookmark string
	}
	MoveBookmarkMsg struct {
		Revision string
		Bookmark string
	}
	DeleteBookmarkMsg struct {
		Bookmark string
	}
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
	SquashOperation
	EditDescriptionOperation
	SetBookmarkOperation
	MoveBookmarkOperation
	DeleteBookmarkOperation
)

type Status int

const (
	Loading Status = iota
	Ready
	Error
)

func Close() tea.Msg {
	return CloseViewMsg{}
}

func Refresh(selectedRevision string) tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg{SelectedRevision: selectedRevision}
	}
}

func UpdateRevSet(revset string) tea.Cmd {
	return func() tea.Msg {
		return UpdateRevSetMsg(revset)
	}
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

func ShowOutput(c *exec.Cmd) tea.Cmd {
	command := strings.Join(c.Args, " ")
	return tea.Batch(CommandRunning(command),
		func() tea.Msg {
			output, err := c.CombinedOutput()
			return CommandCompletedMsg{
				Output: string(output),
				Err:    err,
			}
		})
}
