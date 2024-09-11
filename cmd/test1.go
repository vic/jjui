package main

import (
	"jjui/internal/dag"
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

func main() {
	//commits := TEST_MULTIPLE_CHILDREN
	commits := getJJCommits()
	tree := dag.NewDag()
	for _, commit := range commits {
		tree.AddNode(&commit)
	}

	for _, commit := range commits {
		node := tree.GetNode(&commit)
		for _, parent := range commit.Parents {
			if parent := tree.GetNodeByChangeId(parent); parent != nil {
				parent.AddEdge(node, dag.DirectEdge)
			}
		}
	}

	roots := tree.GetRoots()
	for i := 0; i < len(roots)-1; i++ {
		next := roots[i+1]
		next.AddEdge(roots[i], dag.IndirectEdge)
	}
	dag.Walk(roots[len(roots)-1], dag.DefaultRenderer, dag.RenderContext{Level: 0, IsFirstChild: true})
}

func getJJCommits() []jj.Commit {
	location := os.Getenv("PWD")
	if len(os.Args) > 1 {
		location = os.Args[1]
	}
	return jj.GetCommits(location)
}
