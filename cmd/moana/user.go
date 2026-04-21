package main

import (
	"fmt"
	"os"
)

// userUsage is printed when "moana user" is run without a subcommand.
const userUsage = "usage: moana user <add|password> [flags]\n"

// cliEmailPasswordRequired is printed when user add/password is missing -email or -password.
const cliEmailPasswordRequired = "-email and -password are required\n"

func runUser(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, userUsage)
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
