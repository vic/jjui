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
	SourceDescendants
)

type Target int

const (
	TargetDestination Target = iota
	TargetAfter
	TargetBefore
)

var (
	sourceToFlags = map[Source]string{
		SourceBranch:      "--branch",
		SourceRevision:    "--revisions",
		SourceDescendants: "--source",
	}
	targetToFlags = map[Target]string{
		TargetAfter:       "--insert-after",
		TargetBefore:      "--insert-before",
		TargetDestination: "--destination",
	}
)

type Operation struct {
	context common.AppContext
	From    string
	To      *jj.Commit
	Source  Source
	Target  Target
}

var (
	Revision    = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "change source to revision"))
	Branch      = key.NewBinding(key.WithKeys("B"), key.WithHelp("B", "change source to branch"))
	SourceKey   = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "change source to descendants"))
	Destination = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "change target to destination"))
	After       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "change target to after"))
	Before      = key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "change target to before"))
	Apply       = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "apply"))
	Cancel      = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
)

func (r *Operation) HandleKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, Revision):
		r.Source = SourceRevision
	case key.Matches(msg, Branch):
		r.Source = SourceBranch
	case key.Matches(msg, SourceKey):
		r.Source = SourceDescendants
	case key.Matches(msg, Destination):
		r.Target = TargetDestination
	case key.Matches(msg, After):
		r.Target = TargetAfter
	case key.Matches(msg, Before):
		r.Target = TargetBefore
	case key.Matches(msg, Apply):
		source := sourceToFlags[r.Source]
		target := targetToFlags[r.Target]
		return r.context.RunCommand(jj.Rebase(r.From, r.To.ChangeIdShort, source, target), common.RefreshAndSelect(r.From), common.Close)
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
		SourceKey,
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
	if r.Source == SourceDescendants {
		source = "itself and descendants of "
	}
	if r.Source == SourceRevision {
		source = "only "
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

func NewOperation(context common.AppContext, from string, source Source, target Target) *Operation {
	return &Operation{
		context: context,
		From:    from,
		Source:  source,
		Target:  target,
	}
}
