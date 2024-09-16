package main

import (
	"jjui/internal/jj"
	"os"
)

var TEST = []jj.Commit{
	{ChangeId: "zmuo", Parents: []string{"zurt"}, Description: "(no description)"},
	{ChangeId: "zzty", Parents: []string{"oyyz"}, Description: "(no description)"},
	{ChangeId: "oyyz", Parents: []string{"rkvo"}, Description: "op diff: set up diff/template"},
	{ChangeId: "zurt", Parents: []string{"kvzq"}, Description: "(no description)"},
	{ChangeId: "kvzq", Parents: []string{"xknr"}, Description: "cli_util: add missing word in conflict resolution"},
	{ChangeId: "zovy", Parents: []string{"orrk"}, Description: "readme: include link to wiki"},
	{ChangeId: "orrk", Parents: []string{"lmkm"}, Description: "cli_util: short-prefixes for commit summary"},
}
var TEST2 = []jj.Commit{
	{ChangeId: "top", Parents: []string{"middle"}, Description: "top commit"},
	{ChangeId: "side1", Parents: []string{"side2"}, Description: "side commit"},
	{ChangeId: "side2", Parents: []string{"xyz"}, Description: "side bottom"},
	{ChangeId: "middle", Parents: []string{"bottom"}, Description: "middle commit"},
	{ChangeId: "bottom", Parents: nil, Description: "bottom commit"},
}

var TEST_MULTIPLE_CHILDREN = []jj.Commit{
	{ChangeId: "tymp", Parents: []string{"orrk"}, Description: "top child"},
	{ChangeId: "kvzq", Parents: []string{"orrk"}, Description: "second child"},
	{ChangeId: "zovy", Parents: []string{"orrk"}, Description: "third child"},
	{ChangeId: "orrk", Parents: []string{"lmkm"}, Description: "root commit"},
}

//func main() {
//	//commits := TEST_MULTIPLE_CHILDREN
//	commits := getJJCommits()
//	root := dag.Build(commits)
//	rows := dag.BuildGraphRows(root)
//	builder := strings.Builder{}
//	for _, row := range rows {
//		dag.DefaultRenderer(&builder, &row, dag.DefaultPalette)
//	}
//	fmt.Println(builder.String())
//}

func getJJCommits() []jj.Commit {
	location := os.Getenv("PWD")
	if len(os.Args) > 1 {
		location = os.Args[1]
	}
	return jj.GetCommits(location)
}
