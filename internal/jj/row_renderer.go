package jj

type RowRenderer interface {
	RenderBefore(commit *Commit) string
	RenderAfter(commit *Commit) string
	RenderGlyph(connection ConnectionType, commit *Commit) string
	RenderTermination(connection ConnectionType) string
	RenderChangeId(commit *Commit) string
	RenderCommitId(commit *Commit) string
	RenderAuthor(commit *Commit) string
	RenderDate(commit *Commit) string
	RenderBookmarks(commit *Commit) string
	RenderDescription(commit *Commit) string
	RenderMarkers(commit *Commit) string
}
