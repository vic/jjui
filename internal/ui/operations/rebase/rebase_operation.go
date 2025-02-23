package rebase

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Source int

const (
	SourceRevision Source = iota
	SourceBranch
)

type Target int

const (
	TargetDestination Target = iota
	TargetAfter
	TargetBefore
)

var (
	sourceToFlags = map[Source]string{
		SourceBranch:   "-b",
		SourceRevision: "-r",
	}
	targetToFlags = map[Target]string{
		TargetAfter:       "-A",
		TargetBefore:      "-B",
		TargetDestination: "-d",
	}
)

type Operation struct {
	From     string
	To       *jj.Commit
	Source   Source
	Target   Target
	Commands common.UICommands
}

var (
	Revision    = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rebase revision"))
	Branch      = key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "rebase branch"))
	Destination = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "destination"))
	After       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "after"))
	Before      = key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "before"))
	Apply       = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply"))
	Cancel      = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (r *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Revision):
		r.Source = SourceRevision
	case key.Matches(msg, Branch):
		r.Source = SourceBranch
	case key.Matches(msg, Destination):
		r.Target = TargetDestination
	case key.Matches(msg, After):
		r.Target = TargetAfter
	case key.Matches(msg, Before):
		r.Target = TargetBefore
	case key.Matches(msg, Apply):
		source := sourceToFlags[r.Source]
		target := targetToFlags[r.Target]
		return tea.Batch(r.Commands.Rebase(r.From, r.To.ChangeIdShort, source, target), common.Close)
	case key.Matches(msg, Cancel):
		return common.Close
	}
	return nil
}

func (r *Operation) SetSelectedRevision(commit *jj.Commit) {
	r.To = commit
}

func (r *Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		Revision,
		Branch,
		Destination,
		After,
		Before,
	}
}

func (r *Operation) FullHelp() [][]key.Binding {
	return [][]key.Binding{r.ShortHelp()}
}

func (r *Operation) RenderPosition() operations.RenderPosition {
	if r.Target == TargetAfter {
		return operations.RenderPositionBefore
	}
	if r.Target == TargetDestination {
		return operations.RenderPositionBefore
	}
	return operations.RenderPositionAfter
}

func (r *Operation) Render() string {
	var ret string
	if r.Target == TargetDestination {
		ret = "onto"
	}
	if r.Target == TargetAfter {
		ret = "after"
	}
	if r.Target == TargetBefore {
		ret = "before"
	}
	var source string
	if r.Source == SourceBranch {
		source = "branch of "
	}
	if r.Source == SourceRevision {
		source = ""
	}
	lipgloss.NewStyle().SetString("rebase")
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		common.DropStyle.Render("<< "+ret+" >>"),
		" ",
		common.DefaultPalette.CommitIdRestStyle.Render("rebase"),
		" ",
		common.DefaultPalette.CommitIdRestStyle.Render(source),
		common.DefaultPalette.CommitShortStyle.Render(r.From),
		" ",
		common.DefaultPalette.CommitIdRestStyle.Render(ret),
		" ",
		common.DefaultPalette.CommitShortStyle.Render(r.To.ChangeIdShort),
	)
}

func NewOperation(commands common.UICommands, from string, source Source, target Target) *Operation {
	return &Operation{
		Commands: commands,
		From:     from,
		Source:   source,
		Target:   target,
	}
}
