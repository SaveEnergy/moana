package config

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// minProductionSessionSecretLen is the minimum MOANA_SESSION_SECRET length when
// MOANA_ENV=production (sufficient entropy for HMAC session signing).
const minProductionSessionSecretLen = 32

const (
	defaultSessionMaxAgeSec  = 604800 // 7 days
	defaultRequestTimeoutSec = 60
	defaultPasswordResetMin  = 60
	// defaultRateLimitLoginPerMin and defaultRateLimitForgotPerMin cap abuse of auth POST routes; 0 disables.
	defaultRateLimitLoginPerMin       = 20
	defaultRateLimitForgotPasswordMin = 10
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
	// MaxRequestBodyBytes, if positive, caps HTTP POST body size (MOANA_MAX_REQUEST_BODY_BYTES).
	// Zero means use the server's default limit (1 MiB in [moana/internal/server]).
	MaxRequestBodyBytes int64
	// AvatarDataDir, if set, is the directory for profile JPEGs (MOANA_AVATAR_DIR). Empty means
	// a directory next to the SQLite file ([filepath.Dir] on DBPath + "/avatars") or, for :memory: DB, a host temp path resolved in [moana/internal/app.New].
	AvatarDataDir string
	// PublicBaseURL is the public site origin (e.g. https://moana.example.com) used in password-reset
	// emails. Required when [SendGridAPIKey] is set (MOANA_PUBLIC_BASE_URL).
	PublicBaseURL string
	// PasswordResetTTL is how long a reset link remains valid (MOANA_PASSWORD_RESET_TTL_MIN).
	PasswordResetTTL time.Duration
	// SendGridAPIKey enables outbound mail via SendGrid (MOANA_SENDGRID_API_KEY). Do not log.
	SendGridAPIKey string
	// MailFrom is the verified sender address in SendGrid (MOANA_MAIL_FROM).
	MailFrom string
	// SendGridPasswordResetTemplateID is a SendGrid dynamic template id (d-…) for password reset
	// (MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID). The template should use {{reset_url}} in the body.
	SendGridPasswordResetTemplateID string
	// TrustForwardedAddr treats the first X-Forwarded-For as the client IP (MOANA_TRUST_X_FORWARDED_FOR).
	// Use only when Moana is behind a trusted reverse proxy, not exposed directly to clients.
	TrustForwardedAddr bool
	// RateLimitLoginPerMin limits POST /login per client IP per rolling minute; 0 disables.
	RateLimitLoginPerMin int
	// RateLimitForgotPasswordPerMin limits POST /forgot-password per client IP per rolling minute; 0 disables.
	RateLimitForgotPasswordPerMin int
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

	publicBase := strings.TrimSpace(os.Getenv("MOANA_PUBLIC_BASE_URL"))
	if err := validateLegacySMTPEnv(); err != nil {
		return nil, err
	}
	sendgridKey := strings.TrimSpace(os.Getenv("MOANA_SENDGRID_API_KEY"))
	mailFrom := strings.TrimSpace(os.Getenv("MOANA_MAIL_FROM"))
	sendgridTmpl := strings.TrimSpace(os.Getenv("MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID"))
	if err := validatePublicBaseURLOptional(publicBase); err != nil {
		return nil, err
	}
	if sendgridKey != "" {
		if mailFrom == "" {
			return nil, fmt.Errorf("MOANA_MAIL_FROM is required when MOANA_SENDGRID_API_KEY is set (verified sender email in SendGrid)")
		}
		if publicBase == "" {
			return nil, fmt.Errorf("MOANA_PUBLIC_BASE_URL is required when MOANA_SENDGRID_API_KEY is set (password reset links must be absolute)")
		}
		if sendgridTmpl == "" {
			return nil, fmt.Errorf("MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID is required when MOANA_SENDGRID_API_KEY is set (use a dynamic transactional template; pass {{reset_url}} in the template body)")
		}
	} else {
		if mailFrom != "" {
			return nil, fmt.Errorf("MOANA_MAIL_FROM is set but MOANA_SENDGRID_API_KEY is empty; set both or clear MOANA_MAIL_FROM")
		}
		if sendgridTmpl != "" {
			return nil, fmt.Errorf("MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID is set but MOANA_SENDGRID_API_KEY is empty; set both or clear the template id")
		}
	}

	resetMin := parsePositiveIntEnv("MOANA_PASSWORD_RESET_TTL_MIN", defaultPasswordResetMin)

	avatarDataDir := strings.TrimSpace(os.Getenv("MOANA_AVATAR_DIR"))
	trustFwd := parseBoolTruthy("MOANA_TRUST_X_FORWARDED_FOR")
	rateLogin := parseRateLimitPerMinute("MOANA_RATE_LIMIT_LOGIN_PER_MIN", defaultRateLimitLoginPerMin)
	rateForgot := parseRateLimitPerMinute("MOANA_RATE_LIMIT_FORGOT_PASSWORD_PER_MIN", defaultRateLimitForgotPasswordMin)

	return &Config{
		Listen:                          listen,
		DBPath:                          dbPath,
		AvatarDataDir:                   avatarDataDir,
		SessionSecret:                   secret,
		SecureCookies:                   env == "production",
		SessionMaxAge:                   durationSecondsClamped(maxAgeSec),
		RequestTimeout:                  durationSecondsClamped(timeoutSec),
		RepoURL:                         repoURL,
		MaxRequestBodyBytes:             parseMaxRequestBodyBytesEnv(),
		PublicBaseURL:                   publicBase,
		PasswordResetTTL:                time.Duration(resetMin) * time.Minute,
		SendGridAPIKey:                  sendgridKey,
		MailFrom:                        mailFrom,
		SendGridPasswordResetTemplateID: sendgridTmpl,
		TrustForwardedAddr:              trustFwd,
		RateLimitLoginPerMin:              rateLogin,
		RateLimitForgotPasswordPerMin:     rateForgot,
	}, nil
}

// validateLegacySMTPEnv returns an error if deprecated MOANA_SMTP_* variables are set.
func validateLegacySMTPEnv() error {
	legacy := []struct {
		key, val string
	}{
		{"MOANA_SMTP_HOST", os.Getenv("MOANA_SMTP_HOST")},
		{"MOANA_SMTP_PORT", os.Getenv("MOANA_SMTP_PORT")},
		{"MOANA_SMTP_USER", os.Getenv("MOANA_SMTP_USER")},
		{"MOANA_SMTP_PASSWORD", os.Getenv("MOANA_SMTP_PASSWORD")},
		{"MOANA_SMTP_FROM", os.Getenv("MOANA_SMTP_FROM")},
	}
	for _, row := range legacy {
		if strings.TrimSpace(row.val) == "" {
			continue
		}
		return fmt.Errorf("%s is set but Moana no longer uses SMTP: use MOANA_SENDGRID_API_KEY and MOANA_MAIL_FROM instead (remove or unset legacy MOANA_SMTP_* variables)", row.key)
	}
	return nil
}

func validatePublicBaseURLOptional(raw string) error {
	if raw == "" {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("MOANA_PUBLIC_BASE_URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("MOANA_PUBLIC_BASE_URL must use http or https (got scheme %q)", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("MOANA_PUBLIC_BASE_URL must include a host (e.g. https://moana.example.com)")
	}
	return nil
}

// parseMaxRequestBodyBytesEnv reads MOANA_MAX_REQUEST_BODY_BYTES.
// Empty, invalid, negative, or zero values yield 0 on [Config], so [moana/internal/server.maxRequestBodyBytes] applies the default 1 MiB cap (explicit "0" is not unlimited).
func parseMaxRequestBodyBytesEnv() int64 {
	s := os.Getenv("MOANA_MAX_REQUEST_BODY_BYTES")
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return 0
	}
	return v
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
