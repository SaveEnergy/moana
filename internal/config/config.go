package config

import (
	"fmt"
	"math"
	"os"
	"time"
)

// minProductionSessionSecretLen is the minimum MOANA_SESSION_SECRET length when
// MOANA_ENV=production (sufficient entropy for HMAC session signing).
const minProductionSessionSecretLen = 32

const (
	defaultSessionMaxAgeSec   = 604800 // 7 days
	defaultRequestTimeoutSec = 60
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
	if env == "production" && len(secretStr) < minProductionSessionSecretLen {
		return nil, fmt.Errorf("MOANA_SESSION_SECRET must be at least %d characters when MOANA_ENV=production", minProductionSessionSecretLen)
	}
	var secret []byte
	if secretStr != "" {
		secret = []byte(secretStr)
	} else {
		// Dev-only fallback; not for production.
		secret = []byte("dev-insecure-session-secret-change-me")
	}

	maxAgeSec := parsePositiveIntEnv("MOANA_SESSION_MAX_AGE_SEC", defaultSessionMaxAgeSec)
	timeoutSec := parsePositiveIntEnv("MOANA_REQUEST_TIMEOUT_SEC", defaultRequestTimeoutSec)

	repoURL := getenv("MOANA_REPO_URL", "https://github.com/SaveEnergy/moana")

	return &Config{
		Listen:         listen,
		DBPath:         dbPath,
		SessionSecret:  secret,
		SecureCookies:  env == "production",
		SessionMaxAge:  durationSecondsClamped(maxAgeSec),
		RequestTimeout: durationSecondsClamped(timeoutSec),
		RepoURL:        repoURL,
	}, nil
}

// durationSecondsClamped converts a positive second count to [time.Duration] without int64 overflow.
// Unclamped, [time.Duration](sec)*[time.Second] wraps for very large sec (e.g. negative duration).
func durationSecondsClamped(sec int) time.Duration {
	maxSec := int(math.MaxInt64 / int64(time.Second))
	if sec > maxSec {
		sec = maxSec
	}
	return time.Duration(sec) * time.Second
}
