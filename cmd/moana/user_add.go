package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"moana/internal/auth"
	"moana/internal/store"
)

func runUserAdd(args []string) int {
	fs := flag.NewFlagSet("user add", flag.ExitOnError)
	email := fs.String("email", "", "user email (required)")
	password := fs.String("password", "", "initial password (required)")
	role := fs.String("role", "user", "user or admin")
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

	st, sqlDB, err := openCLIStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		return 1
	}
	defer sqlDB.Close()

	hash, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hash: %v\n", err)
		return 1
	}
	ctx := context.Background()
	id, err := st.CreateUser(ctx, *email, hash, r)
	if err != nil {
		if errors.Is(err, store.ErrDuplicateUserEmail) {
			fmt.Fprintf(os.Stderr, "create user: email already exists\n")
			return 1
		}
		if errors.Is(err, store.ErrInvalidUserEmail) {
			fmt.Fprintf(os.Stderr, "create user: invalid email\n")
			return 1
		}
		fmt.Fprintf(os.Stderr, "create user: %v\n", err)
		return 1
	}
	fmt.Printf("created user id=%d email=%s role=%s\n", id, strings.TrimSpace(*email), r)
	return 0
}
