package dag

import (
	"testing"

	"jjui/internal/jj"

	"github.com/stretchr/testify/assert"
)

func TestBuildGraphRows_WithElidedCommits(t *testing.T) {
	commits := []jj.Commit{
		{
			ChangeId: "topchange",
			Parents:  nil,
		},
		{
			ChangeId: "psvvky",
			Parents:  []string{"zzzzzz"},
		},
		{
			ChangeId: "zzzzzz",
			Parents:  nil,
		},
	}
	root := Build(commits)
	rows := BuildGraphRows(root)
	assert.Len(t, rows, 3)
	assert.Equal(t, commits[0].ChangeId, rows[0].Commit.ChangeId)
	assert.Equal(t, commits[1].ChangeId, rows[1].Commit.ChangeId)
	assert.Equal(t, commits[2].ChangeId, rows[2].Commit.ChangeId)
}

func TestBuildGraphRows_WithTop2Commits(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "top_empty", Parents: []string{"parent"}},
		{ChangeId: "top_addfile", Parents: []string{"parent"}},
		{ChangeId: "parent", Parents: nil},
	}
	root := Build(commits)
	rows := BuildGraphRows(root)

	assert.Len(t, rows, len(commits))
	sortedChangeIds := []string{rows[0].Commit.ChangeId, rows[1].Commit.ChangeId, rows[2].Commit.ChangeId}
	assert.Exactly(t, []string{"top_empty", "top_addfile", "parent"}, sortedChangeIds)
	assert.Equal(t, 1, rows[1].Level)
}

func TestBuildGraphRows_LevelsWithElidedRevisions(t *testing.T) {
	commits := []jj.Commit{
		{ChangeId: "top", Parents: nil},
		{ChangeId: "middle", Parents: []string{"middle_parent"}},
		{ChangeId: "middle_parent", Parents: nil},
	}
	root := Build(commits)
	rows := BuildGraphRows(root)
	assert.Len(t, rows, len(commits))
	assert.Equal(t, 0, rows[0].Level, "top should be at level 0")
	assert.Equal(t, 1, rows[1].Level, "middle should be at level 1")
	assert.Equal(t, 0, rows[2].Level, "middle_parent should be at level 0")
}
