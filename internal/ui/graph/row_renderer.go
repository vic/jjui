package graph

import "github.com/idursun/jjui/internal/jj"

type RowRenderer interface {
	RenderBefore(commit *jj.Commit) string
	RenderAfter(commit *jj.Commit) string
	RenderGlyph(connection jj.ConnectionType, commit *jj.Commit) string
	RenderTermination(connection jj.ConnectionType) string
	RenderChangeId(commit *jj.Commit) string
	RenderCommitId(commit *jj.Commit) string
	RenderAuthor(commit *jj.Commit) string
	RenderDate(commit *jj.Commit) string
	RenderBookmarks(commit *jj.Commit) string
	RenderDescription(commit *jj.Commit) string
	RenderMarkers(commit *jj.Commit) string
}
