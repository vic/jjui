package rebase

import (
	"github.com/charmbracelet/bubbles/key"
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

type Operation struct {
	From   string
	To     *jj.Commit
	Source Source
	Target Target
}

func (r *Operation) SetSelectedRevision(commit *jj.Commit) {
	r.To = commit
}

func (r *Operation) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
	}
}

func (r *Operation) FullHelp() [][]key.Binding {
	//TODO implement me
	panic("implement me")
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
	if r.Target == TargetDestination {
		return common.DropStyle.Render("<< onto >> " + r.To.ChangeIdShort)
	}
	if r.Target == TargetAfter {
		return common.DropStyle.Render("<< after >>")
	}
	if r.Target == TargetBefore {
		return common.DropStyle.Render("<< before >>")
	}
	return ""
}

func (r *Operation) GetSourceTargetFlags() (source string, target string) {
	switch r.Source {
	case SourceBranch:
		source = "-b"
	case SourceRevision:
		source = "-r"
	}
	switch r.Target {
	case TargetAfter:
		target = "-A"
	case TargetBefore:
		target = "-B"
	case TargetDestination:
		target = "-d"
	}
	return source, target
}

func NewOperation(from string, source Source, target Target) *Operation {
	return &Operation{
		From:   from,
		Source: source,
		Target: target,
	}
}
