package store

import "time"

// User is an account that can sign in.
type User struct {
	ID            int64
	Email         string
	PasswordHash  []byte
	Role          string
	CreatedAt     time.Time
	HouseholdID   int64
	FirstName     string
	LastName      string
	HouseholdRole string
}

// UserSummary is a user row without secret fields (for admin lists).
type UserSummary struct {
	ID        int64
	Email     string
	Role      string
	CreatedAt time.Time
}

// HouseholdMember is a household row for display (no secrets).
type HouseholdMember struct {
	ID            int64
	Email         string
	FirstName     string
	LastName      string
	HouseholdRole string
}
