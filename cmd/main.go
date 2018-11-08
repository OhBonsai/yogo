package main

import (
	"github.com/OhBonsai/yogo/cmd/commands"
	"os"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}