package main

import (
	"fmt"
	"jjui/internal/jj"
	"os"
)

func main() {
	//commits := jj.GetCommits("/Users/idursun/repositories/elixir/beach_games")
	commits := jj.GetCommits(os.Getenv("PWD"))
	for _, commit := range commits {
		fmt.Printf("%+v\n", commit)
	}
}
