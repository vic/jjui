package helppage

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
	"github.com/idursun/jjui/internal/ui/operations/bookmark"
	"github.com/idursun/jjui/internal/ui/operations/git"
	"github.com/idursun/jjui/internal/ui/operations/rebase"
	"github.com/idursun/jjui/internal/ui/operations/squash"
)

type Model struct {
	width  int
	height int
}

func (h *Model) Width() int {
	return h.width
}

func (h *Model) Height() int {
	return h.height
}

func (h *Model) SetWidth(w int) {
	h.width = w
}

func (h *Model) SetHeight(height int) {
	h.height = height
}

func (h *Model) ShortHelp() []key.Binding {
	return []key.Binding{operations.Help, operations.Cancel}
}

func (h *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{h.ShortHelp()}
}

func (h *Model) Init() tea.Cmd {
	return nil
}

func (h *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, operations.Help), key.Matches(msg, operations.Cancel):
			return h, common.Close
		}
	}
	return h, nil
}

var (
	keyStyle  = common.DefaultPalette.CommitShortStyle
	descStyle = common.DefaultPalette.CommitIdRestStyle
)

func printHelp(k key.Binding) string {
	return printHelpExt(k.Help().Key, k.Help().Desc)
}

func printHelpExt(key string, desc string) string {
	keyAligned := fmt.Sprintf("%9s", key)
	help := fmt.Sprintf("%s %s", keyStyle.Render(keyAligned), descStyle.Render(desc))
	return help
}

func printHeader(header string) string {
	return common.DefaultPalette.Empty.Render(header)
}

func printMode(key key.Binding, name string) string {
	keyAligned := fmt.Sprintf("%9s", key.Help().Key)
	help := fmt.Sprintf("%v %s", keyStyle.Render(keyAligned), common.DefaultPalette.Empty.Render(name))
	return help
}

func (h *Model) View() string {
	leftView := lipgloss.JoinVertical(lipgloss.Left,
		printHeader("UI"),
		printHelp(operations.Refresh),
		printHelp(operations.Help),
		printHelp(operations.Cancel),
		printHelp(operations.Quit),
		"",
		printHeader("Preview"),
		printHelpExt("tab", "focus preview"),
		printHelpExt("ctrl+p", "line up"),
		printHelpExt("ctrl+n", "line down"),
		printHelpExt("ctrl+d", "half page down"),
		printHelpExt("ctrl+u", "half page up"),
		"",
		printHeader("Preview (when focused)"),
		printHelpExt("tab", "unfocus preview"),
		printHelpExt("k", "line up"),
		printHelpExt("j", "line down"),
		printHelpExt("d", "half page down"),
		printHelpExt("u", "half page up"),
		"",
		printHeader("Revisions"),
		printHelp(operations.New),
		printHelp(operations.Description),
		printHelp(operations.Edit),
		printHelp(operations.Diff),
		printHelp(operations.Diffedit),
		printHelp(operations.Split),
		printHelp(operations.Abandon),
		printHelp(operations.Undo),
		"",
		printHeader("Revset"),
		printHelp(operations.Revset),
	)

	rightView := lipgloss.JoinVertical(lipgloss.Left,
		printMode(operations.GitMode, "Git"),
		printHelp(git.Push),
		printHelp(git.Fetch),
		"",
		printMode(operations.BookmarkMode, "Bookmarks"),
		printHelp(bookmark.Move),
		printHelp(bookmark.Set),
		printHelp(bookmark.Delete),
		"",
		printMode(operations.RebaseMode, "Rebase"),
		printHelp(rebase.Before),
		printHelp(rebase.After),
		printHelp(rebase.Destination),
		printHelp(rebase.Revision),
		printHelp(rebase.SourceKey),
		printHelp(rebase.Apply),
		"",
		printMode(operations.SquashMode, "Squash"),
		printHelp(squash.Apply),
	)

	content := lipgloss.JoinHorizontal(lipgloss.Left, leftView, "  ", rightView)

	return lipgloss.Place(h.width, h.height, lipgloss.Center, lipgloss.Center, content)
}
