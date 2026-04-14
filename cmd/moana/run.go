package main

import (
	"fmt"
	"os"
)

func run(args []string) int {
	if len(args) >= 2 && args[1] == "user" {
		return runUser(args[2:])
	}
	if len(args) >= 2 && args[1] == "serve" {
		runServe()
		return 0
	}
	if len(args) >= 2 {
		fmt.Fprintf(os.Stderr, "usage: moana [serve | user <add|password> ...]\n")
		return 1
	}
	runServe()
	return 0
}
