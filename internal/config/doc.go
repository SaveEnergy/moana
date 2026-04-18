// Package config loads [Config] from environment variables used by cmd/moana and tests.
//
// Keys include MOANA_LISTEN, MOANA_DB_PATH, MOANA_ENV, MOANA_SESSION_SECRET,
// MOANA_SESSION_MAX_AGE_SEC (default 7d), MOANA_REQUEST_TIMEOUT_SEC (default 60s), and MOANA_REPO_URL — see [Load].
// Regression coverage lives in config_test.go.
package config
