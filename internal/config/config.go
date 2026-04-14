package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds runtime settings loaded from the environment.
type Config struct {
	Listen         string
	DBPath         string
	SessionSecret  []byte
	SecureCookies  bool
	SessionMaxAge  time.Duration
	RequestTimeout time.Duration
	// RepoURL is the public source repository (e.g. GitHub), shown in the app footer.
	RepoURL string
}

// Load reads configuration from the environment. MOANA_SESSION_SECRET is required
// when MOANA_ENV is production.
func Load() (*Config, error) {
	listen := getenv("MOANA_LISTEN", ":8080")
	dbPath := getenv("MOANA_DB_PATH", "data/moana.db")
	env := getenv("MOANA_ENV", "development")

	secretStr := os.Getenv("MOANA_SESSION_SECRET")
	if env == "production" && secretStr == "" {
		return nil, fmt.Errorf("MOANA_SESSION_SECRET is required when MOANA_ENV=production")
	}
	var secret []byte
	if secretStr != "" {
		secret = []byte(secretStr)
	} else {
		// Dev-only fallback; not for production.
		secret = []byte("dev-insecure-session-secret-change-me")
	}

	maxAgeSec, _ := strconv.Atoi(getenv("MOANA_SESSION_MAX_AGE_SEC", "604800"))
	if maxAgeSec <= 0 {
		maxAgeSec = 604800
	}
	timeoutSec, _ := strconv.Atoi(getenv("MOANA_REQUEST_TIMEOUT_SEC", "60"))
	if timeoutSec <= 0 {
		timeoutSec = 60
	}

	repoURL := getenv("MOANA_REPO_URL", "https://github.com/sinan/moana")

	return &Config{
		Listen:         listen,
		DBPath:         dbPath,
		SessionSecret:  secret,
		SecureCookies:  env == "production",
		SessionMaxAge:  time.Duration(maxAgeSec) * time.Second,
		RequestTimeout: time.Duration(timeoutSec) * time.Second,
		RepoURL:        repoURL,
	}, nil
}

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
