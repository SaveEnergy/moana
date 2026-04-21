package main

import (
	"fmt"
	"os"
)

// runUsage is printed for unknown top-level commands and invalid "serve" invocation.
const runUsage = "usage: moana [serve | user <add|password> ...]\n"

func run(args []string) int {
	if len(args) >= 2 && args[1] == "user" {
		return runUser(args[2:])
	}
	if len(args) >= 2 && args[1] == "serve" {
		if len(args) > 2 {
			fmt.Fprint(os.Stderr, runUsage)
			return 1
		}
		runServe()
		return 0
	}
	if len(args) >= 2 {
		fmt.Fprint(os.Stderr, runUsage)
		return 1
	}
	runServe()
	return 0
}
