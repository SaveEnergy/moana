package config

import (
	"os"
	"strconv"
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
