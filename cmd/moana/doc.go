// Command moana is the Moana binary: default and "serve" run the HTTP server; "user" hosts
// account subcommands (add, password). Database path and session config come from the environment (see config.Load).
// run_test.go covers CLI dispatch edge cases (unknown subcommand, user without subcommand).
package main
