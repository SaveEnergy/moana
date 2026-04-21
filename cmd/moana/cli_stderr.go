package main

import (
	"fmt"
	"os"
)

func stderrDatabase(err error) {
	fmt.Fprintf(os.Stderr, "database: %v\n", err)
}

func stderrHash(err error) {
	fmt.Fprintf(os.Stderr, "hash: %v\n", err)
}

func stderrLookup(err error) {
	fmt.Fprintf(os.Stderr, "lookup: %v\n", err)
}

func stderrUpdate(err error) {
	fmt.Fprintf(os.Stderr, "update: %v\n", err)
}

func stderrCreateUser(err error) {
	fmt.Fprintf(os.Stderr, "create user: %v\n", err)
}
