// Command moana is the Moana binary: default and "serve" run the HTTP server; "user" hosts
// account subcommands (add, password). Database path and session config come from the environment (see config.Load).
// run.go defines [runUsage] for consistent stderr copy; user.go defines [userUsage] and [cliEmailPasswordRequired] for the user subcommand tree.
// cli_stderr.go centralizes stderr helpers (database/hash/lookup/update/create-user errors) for user add/password; user.go holds fixed CLI strings including [cliRoleInvalid] and create-user messages.
// run_test.go covers dispatch edge cases (unknown subcommand, user without subcommand, extra args after serve, user add/password without required flags, invalid user add role, user password for unknown email with in-memory DB).
package main
