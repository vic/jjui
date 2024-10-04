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
timestamp
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
		Timestamp:     "timestamp",
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
timestamp
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
		Timestamp:     "timestamp",
		Description:   "add test",
	}

	assert.EqualExportedValues(t, expected, commit)
}
