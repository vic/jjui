package main

import (
	"fmt"
	"jjui/internal/jj"
	"os"
)

func main() {
	// get argument
	location := os.Getenv("PWD")
	if len(os.Args) > 1 {
		location = os.Args[1]
	}
	commits := jj.GetCommits(location)
	for _, commit := range commits {
		fmt.Printf("%+v\n", commit)
	}
}
