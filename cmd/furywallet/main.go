package main

import (
	"os"

	cmd "github.com/elysiumstation/fury/cmd/furywallet/commands"
)

func main() {
	writer := &cmd.Writer{
		Out: os.Stdout,
		Err: os.Stderr,
	}
	cmd.Execute(writer)
}
