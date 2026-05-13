package main

import (
	"fmt"
	"os"

	"github.com/sylphy/git-switch/cli/commands"
)

func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
