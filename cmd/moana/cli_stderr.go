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
