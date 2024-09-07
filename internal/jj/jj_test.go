package jj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseLogOutput_Single(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
ps
psvvky
042d3f0018e5ae891ee2452274b3c7832d33cd5e
ibrahim dursun <some@email.cc>
!!NONE
add test binary

__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "ps",
			ChangeId:      "psvvky",
			CommitId:      "042d3f0018e5ae891ee2452274b3c7832d33cd5e",
			Author:        "ibrahim dursun <some@email.cc>",
			Description:   "add test binary",
		},
	}

	assert.Equal(t, expected, commits)
}

func Test_parseLogOutput_RootCommit(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
z
zzzzzz
0000000000000000000000000000000000000000
!!NONE
!!NONE
__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			CommitId:      "0000000000000000000000000000000000000000",
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
042d3f0018e5ae891ee2452274b3c7832d33cd5e
ibrahim dursun <some@email.cc>
!!NONE
add test binary

__END__
__BEGIN__
z
zzzzzz
0000000000000000000000000000000000000000
!!NONE
!!NONE
__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "ps",
			ChangeId:      "psvvky",
			CommitId:      "042d3f0018e5ae891ee2452274b3c7832d33cd5e",
			Author:        "ibrahim dursun <some@email.cc>",
			Description:   "add test binary",
		},
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			CommitId:      "0000000000000000000000000000000000000000",
			Author:        "",
			Description:   "",
			Branches:      "",
		},
	}

	assert.Equal(t, expected, commits)
}
