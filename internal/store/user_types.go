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
	// AvatarRev is 0 if no custom photo; increment when a new image is stored (served at /avatars/{id} with ?v=).
	AvatarRev int64
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
	AvatarRev     int64
}
