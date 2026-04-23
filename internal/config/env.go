package config

import (
	"os"
	"strconv"
	"strings"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parsePositiveIntEnv reads key as a base-10 integer; on parse error or non-positive value, returns defaultVal.
func parsePositiveIntEnv(key string, defaultVal int) int {
	v, err := strconv.Atoi(getenv(key, strconv.Itoa(defaultVal)))
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}

// DBPath returns MOANA_DB_PATH or the default file path (for CLI tools that only touch the database).
func DBPath() string {
	return getenv("MOANA_DB_PATH", "data/moana.db")
}

// parseBoolTruthy returns true when MOANA_* is 1, true, or yes (case-insensitive).
func parseBoolTruthy(key string) bool {
	s := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return s == "1" || s == "true" || s == "yes"
}

// parseRateLimitPerMinute reads key as a non-negative integer per rolling minute; empty uses defaultVal;
// invalid negative values fall back to defaultVal; explicit 0 disables the limiter.
func parseRateLimitPerMinute(key string, defaultVal int) int {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return defaultVal
	}
	return v
}
