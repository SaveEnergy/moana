package main

import (
	"fmt"
	"os"
)

// version is the binary version string; override at link time, e.g.
// go build -ldflags "-X main.version=1.2.3"
var version = "dev"

// runUsage is printed for unknown top-level commands and invalid "serve" invocation.
const runUsage = "usage: moana [serve | version | help | user <add|password> ...]\n"

func run(args []string) int {
	if len(args) >= 2 && helpArg(args[1]) {
		if len(args) > 2 {
			fmt.Fprint(os.Stderr, runUsage)
			return 1
		}
		fmt.Fprint(os.Stdout, runUsage)
		return 0
	}
	if len(args) >= 2 && versionArg(args[1]) {
		if len(args) > 2 {
			fmt.Fprint(os.Stderr, runUsage)
			return 1
		}
		fmt.Fprintln(os.Stdout, version)
		return 0
	}
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

func versionArg(s string) bool {
	switch s {
	case "version", "-version", "--version":
		return true
	default:
		return false
	}
}

func helpArg(s string) bool {
	switch s {
	case "help", "-h", "--help", "-help":
		return true
	default:
		return false
	}
}
