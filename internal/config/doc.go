// Package config loads [Config] from environment variables used by cmd/moana and tests.
//
// Keys include MOANA_LISTEN, MOANA_DB_PATH, MOANA_ENV, MOANA_SESSION_SECRET,
// MOANA_SESSION_MAX_AGE_SEC (default 7d), MOANA_REQUEST_TIMEOUT_SEC (default 60s), MOANA_REPO_URL,
// optional MOANA_MAX_REQUEST_BODY_BYTES (positive bytes; unset uses the server default 1 MiB cap), and
// optional outbound mail: MOANA_SENDGRID_API_KEY, MOANA_MAIL_FROM,
// MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID, MOANA_PUBLIC_BASE_URL (password reset) — see [Load].
//
// Rate limits (non-negative integers per rolling minute; empty uses server defaults, 0 disables the limiter):
// MOANA_RATE_LIMIT_LOGIN_PER_MIN, MOANA_RATE_LIMIT_FORGOT_PASSWORD_PER_MIN.
// MOANA_TRUST_X_FORWARDED_FOR (1/true/yes) uses the first X-Forwarded-For address as the client for those limits; use only behind a trusted reverse proxy.
//
// Regression coverage lives in config_test.go.
package config
