package config

import "os"

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DBPath returns MOANA_DB_PATH or the default file path (for CLI tools that only touch the database).
func DBPath() string {
	return getenv("MOANA_DB_PATH", "data/moana.db")
}
