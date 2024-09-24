package jj

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_parseLogOutput_Single(t *testing.T) {
	commits := parseLogOutput(`__BEGIN__
m
myxzmzolxmpz
psvvkyllponl
true
false
false
false
ibrahim dursun <some@email.cc>
main
add test
__END__`)

	expected := Commit{
		ChangeIdShort: "m",
		ChangeId:      "myxzmzolxmpz",
		Parents:       []string{"psvvkyllponl"},
		IsWorkingCopy: true,
		Immutable:     false,
		Conflict:      false,
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
false
false
false
!!NONE
!!NONE
__END__`)

	expected := []Commit{
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			IsWorkingCopy: false,
			Immutable:     false,
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
false
false
false
ibrahim dursun <some@email.cc>
!!NONE
add test binary

__END__
__BEGIN__
z
zzzzzz
!!NONE
false
false
false
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
			Index:         0,
		},
		{
			ChangeIdShort: "z",
			ChangeId:      "zzzzzz",
			Parents:       nil,
			IsWorkingCopy: false,
			Author:        "",
			Description:   "",
			Branches:      "",
			Index:         1,
		},
	}

	assert.EqualExportedValues(t, expected[0], commits[0])
	assert.EqualExportedValues(t, expected[1], commits[1])
}

func Test_parseCommit(t *testing.T) {
	lines := strings.Split(`__BEGIN__
c
current
parent
true
false
true
false
ibrahim dursun <some@email.cc>
main
add test
__END__`, "\n")
	commit := parseCommit(lines)

	expected := Commit{
		ChangeIdShort: "c",
		ChangeId:      "current",
		Parents:       []string{"parent"},
		IsWorkingCopy: true,
		Immutable:     false,
		Conflict:      true,
		Author:        "ibrahim dursun <some@email.cc>",
		Branches:      "main",
		Description:   "add test",
	}

	assert.EqualExportedValues(t, expected, commit)
}
