package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"moana/internal/auth"
)

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

	st, sqlDB, err := openCLIStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		return 1
	}
	defer sqlDB.Close()

	ctx := context.Background()
	u, err := st.GetUserByEmail(ctx, *email)
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
