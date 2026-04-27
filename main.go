// main.go is the entry point for the portmap-cli application.
// It delegates execution to the cmd package's root command.
package main

import (
	"os"

	"github.com/user/portmap-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
