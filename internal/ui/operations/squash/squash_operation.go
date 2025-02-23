package squash

import (
	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/operations"
)

type Operation struct {
	From    string
	Current *jj.Commit
}

func (s *Operation) SetSelectedRevision(commit *jj.Commit) {
	s.Current = commit
}

func (s *Operation) Render() string {
	if s.Current == nil {
		return common.DropStyle.Render("<< into >>")
	} else {
		return common.DropStyle.Render(s.Current.ChangeIdShort + " << into >>")
	}
}

func (s *Operation) RenderPosition() operations.RenderPosition {
	return operations.RenderPositionGlyph
}

func NewOperation(from string) *Operation {
	return &Operation{
		From: from,
	}
}
