package server

import (
	"net/http"
	"time"

	"moana/internal/handlers"
)

// larger than [defaultMaxRequestBodyBytes] for multipart image uploads; see [WithMaxRequestBodyBytesAdaptive].
const settingsAvatarPostMaxBytes = 4 * 1 << 20

func requestTimeout(opts *RouterOptions, app *handlers.App) time.Duration {
	if opts != nil && opts.RequestTimeout > 0 {
		return opts.RequestTimeout
	}
	if app != nil && app.Config != nil && app.Config.RequestTimeout > 0 {
		return app.Config.RequestTimeout
	}
	return 0
}

func maxRequestBodyBytes(opts *RouterOptions) int64 {
	if opts == nil || opts.MaxRequestBodyBytes == 0 {
		return defaultMaxRequestBodyBytes
	}
	if opts.MaxRequestBodyBytes < 0 {
		return 0
	}
	return opts.MaxRequestBodyBytes
}

// wrapRouterMiddleware applies POST body size limits and request deadlines around the mux (static,
// health, app routes). Order matches the historical stack: [WithPostAuthRateLimit] is outermost
// (reject abuse before body reads and timeouts), then [WithRequestTimeout], then
// [WithMaxRequestBodyBytesAdaptive] so the deadline covers body reads and handlers. This helper is
// the single composition used by [NewRouterWithRouterOptions] and by tests that assert middleware
// behavior.
func wrapRouterMiddleware(inner http.Handler, opts *RouterOptions, app *handlers.App) http.Handler {
	inner = WithMaxRequestBodyBytesAdaptive(
		maxRequestBodyBytes(opts),
		http.MethodPost,
		handlers.SettingsAvatarPath,
		settingsAvatarPostMaxBytes,
	)(inner)
	if d := requestTimeout(opts, app); d > 0 {
		inner = WithRequestTimeout(d)(inner)
	}
	inner = WithPostAuthRateLimit(app)(inner)
	return inner
}
