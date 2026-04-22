// Package config loads [Config] from environment variables used by cmd/moana and tests.
//
// Keys include MOANA_LISTEN, MOANA_DB_PATH, MOANA_ENV, MOANA_SESSION_SECRET,
// MOANA_SESSION_MAX_AGE_SEC (default 7d), MOANA_REQUEST_TIMEOUT_SEC (default 60s), MOANA_REPO_URL,
// optional MOANA_MAX_REQUEST_BODY_BYTES (positive bytes; unset uses the server default 1 MiB cap), and
// optional outbound mail: MOANA_SENDGRID_API_KEY, MOANA_MAIL_FROM, MOANA_PUBLIC_BASE_URL (password reset) — see [Load].
// Regression coverage lives in config_test.go.
package config
