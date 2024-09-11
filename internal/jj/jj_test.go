package jj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseLogOutput_Single(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
m
myxzmzolxmpz
psvvkyllponl
true
ibrahim dursun <some@email.cc>
main
add test

__END__`)

	expected := Commit{
		ChangeIdShort: "m",
		ChangeId:      "myxzmzolxmpz",
		Parents:       nil,
		IsWorkingCopy: true,
		Author:        "ibrahim dursun <some@email.cc>",
		Branches:      "main",
		Description:   "add test",
	}

	assert.EqualExportedValues(t, expected, commits[0])
}

func Test_parseLogOutput_RootCommit(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
z
zzzzzz
!!NONE
false
!!NONE
!!NONE
__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			IsWorkingCopy: false,
			Author:        "",
			Description:   "",
			Branches:      "",
		},
	}

	assert.Equal(t, expected, commits)
}

func Test_parseLogOutput_TwoCommits(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
ps
psvvky
zzzzzz
true
ibrahim dursun <some@email.cc>
!!NONE
add test binary

__END__
__BEGIN__
z
zzzzzz
!!NONE
false
!!NONE
!!NONE
__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "ps",
			ChangeId:      "psvvky",
			Parents:       []string{"zzzzzz"},
			IsWorkingCopy: true,
			Author:        "ibrahim dursun <some@email.cc>",
			Description:   "add test binary",
		},
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			Parents:       nil,
			IsWorkingCopy: false,
			Author:        "",
			Description:   "",
			Branches:      "",
		},
	}

	assert.EqualExportedValues(t, expected[0], commits[0])
	assert.EqualExportedValues(t, expected[1], commits[1])
}

func TestBuildCommitTree_WithElidedCommits(t *testing.T) {
	commits := []Commit{
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
	sorted := BuildCommitTree(commits)
	assert.Len(t, sorted, 3)
	assert.Equal(t, commits[0].ChangeId, sorted[0].ChangeId)
	assert.Equal(t, commits[1].ChangeId, sorted[1].ChangeId)
	assert.Equal(t, commits[2].ChangeId, sorted[2].ChangeId)
}

func TestBuildCommitTree_WithTop2Commits(t *testing.T) {
	commits := []Commit{
		{ChangeId: "top_empty", Parents: []string{"parent"}},
		{ChangeId: "top_addfile", Parents: []string{"parent"}},
		{ChangeId: "parent", Parents: nil},
	}
	sorted := BuildCommitTree(commits)

	assert.Len(t, sorted, len(commits))
	sortedChangeIds := []string{sorted[0].ChangeId, sorted[1].ChangeId, sorted[2].ChangeId}
	assert.Exactly(t, []string{"top_empty", "top_addfile", "parent"}, sortedChangeIds)
	assert.Equal(t, 1, sorted[1].Level())
}

func TestBuildCommitTree_LevelsWithElidedRevisions(t *testing.T) {
	commits := []Commit{
		{ChangeId: "top", Parents: nil},
		{ChangeId: "middle", Parents: []string{"middle_parent"}},
		{ChangeId: "middle_parent", Parents: nil},
	}
	sorted := BuildCommitTree(commits)
	assert.Len(t, sorted, len(commits))
	assert.Equal(t, 0, sorted[0].Level(), "top should be at level 0")
	assert.Equal(t, 1, sorted[1].Level(), "middle should be at level 1")
	assert.Equal(t, 0, sorted[2].Level(), "middle_parent should be at level 0")
}
