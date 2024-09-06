package main

import (
	"fmt"
	"jjui/internal/jj"
)

func main() {
	commits := jj.GetCommits("/Users/idursun/repositories/elixir/beach_games")
	for _, commit := range commits {
		fmt.Printf("%+v\n", commit)
	}
}
