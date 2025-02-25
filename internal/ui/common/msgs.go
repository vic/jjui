package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/operations"
	"strings"
)

type (
	CloseViewMsg             struct{}
	RefreshMsg               struct{ SelectedRevision string }
	SelectRevisionMsg        string
	SetOperationMsg          struct{ Operation operations.Operation }
	ShowDiffMsg              string
	UpdateRevSetMsg          string
	UpdateRevisionsMsg       []jj.GraphRow
	UpdateRevisionsFailedMsg error
	UpdateBookmarksMsg       struct {
		Bookmarks []string
		Revision  string
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
	UpdatePreviewChangeIdMsg struct {
		ChangeId string
	}
	UpdatePreviewContentMsg struct {
		Content string
	}
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

func SetOperation(op operations.Operation) tea.Cmd {
	return func() tea.Msg {
		return SetOperationMsg{Operation: op}
	}
}

func CommandRunning(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandRunningMsg(command)
	}
}

func RunCommand(c jj.Command, continuations ...tea.Cmd) tea.Cmd {
	commands := make([]tea.Cmd, 0)
	commands = append(commands,
		func() tea.Msg {
			_, err := c.CombinedOutput()
			return CommandCompletedMsg{
				// Content: string(output),
				Err: err,
			}
		})
	commands = append(commands, continuations...)
	return tea.Batch(
		CommandRunning(strings.Join(c.Args(), " ")),
		tea.Sequence(commands...),
	)
}
