package main

import (
	"fmt"
	"os"
)

func runUser(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: moana user <add|password> [flags]\n")
		return 2
	}
	switch args[0] {
	case "add":
		return runUserAdd(args[1:])
	case "password":
		return runUserPassword(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown user subcommand %q\n", args[0])
		return 2
	}
}
