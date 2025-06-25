package graph

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/idursun/jjui/internal/ui/common"

	"github.com/idursun/jjui/internal/jj"
	"github.com/idursun/jjui/internal/ui/operations"
)

type DefaultRowDecorator struct {
	Palette             common.Palette
	HighlightBackground lipgloss.AdaptiveColor
	SearchText          string
	IsHighlighted       bool
	IsSelected          bool
	Op                  operations.Operation
	Width               int
}

func (s *DefaultRowDecorator) RenderBefore(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderPositionBefore)
}

func (s *DefaultRowDecorator) RenderAfter(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderPositionAfter)
}

func (s *DefaultRowDecorator) RenderBeforeChangeId(commit *jj.Commit) string {
	opMarker := s.Op.Render(commit, operations.RenderBeforeChangeId)
	selectedMarker := ""
	if s.IsSelected {
		if s.IsHighlighted {
			selectedMarker = s.Palette.Added.Background(s.HighlightBackground).Render("✓ ")
		} else {
			selectedMarker = s.Palette.Added.Render("✓ ")
		}
	}
	return opMarker + selectedMarker
}

func (s *DefaultRowDecorator) RenderBeforeCommitId(commit *jj.Commit) string {
	return s.Op.Render(commit, operations.RenderBeforeCommitId)
}
