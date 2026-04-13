package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"moana/internal/auth"
	"moana/internal/config"
	"moana/internal/db"
	"moana/internal/store"
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

func runUserAdd(args []string) int {
	fs := flag.NewFlagSet("user add", flag.ExitOnError)
	email := fs.String("email", "", "user email (required)")
	password := fs.String("password", "", "initial password (required)")
	role := fs.String("role", "user", "user or admin")
	tz := fs.String("timezone", "UTC", "IANA time zone for display (e.g. Europe/Berlin)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if strings.TrimSpace(*email) == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "-email and -password are required\n")
		return 2
	}
	r := strings.ToLower(strings.TrimSpace(*role))
	if r != "user" && r != "admin" {
		fmt.Fprintf(os.Stderr, "-role must be user or admin\n")
		return 2
	}

	database, err := db.Open(config.DBPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		return 1
	}
	defer database.Close()

	hash, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hash: %v\n", err)
		return 1
	}

	st := store.New(database)
	ctx := context.Background()
	id, err := st.CreateUser(ctx, strings.TrimSpace(*email), hash, r, strings.TrimSpace(*tz))
	if err != nil {
		fmt.Fprintf(os.Stderr, "create user: %v\n", err)
		return 1
	}
	fmt.Printf("created user id=%d email=%s role=%s\n", id, strings.TrimSpace(*email), r)
	return 0
}

func runUserPassword(args []string) int {
	fs := flag.NewFlagSet("user password", flag.ExitOnError)
	email := fs.String("email", "", "user email (required)")
	password := fs.String("password", "", "new password (required)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if strings.TrimSpace(*email) == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "-email and -password are required\n")
		return 2
	}

	database, err := db.Open(config.DBPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		return 1
	}
	defer database.Close()

	st := store.New(database)
	ctx := context.Background()
	u, err := st.GetUserByEmail(ctx, strings.TrimSpace(*email))
	if err != nil {
		fmt.Fprintf(os.Stderr, "lookup: %v\n", err)
		return 1
	}
	if u == nil {
		fmt.Fprintf(os.Stderr, "user not found\n")
		return 1
	}
	hash, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hash: %v\n", err)
		return 1
	}
	if err := st.UpdateUserPassword(ctx, u.ID, hash); err != nil {
		fmt.Fprintf(os.Stderr, "update: %v\n", err)
		return 1
	}
	fmt.Printf("updated password for %s\n", strings.TrimSpace(*email))
	return 0
}
