package handlers

// Route path constants for redirects. Keep in sync with routes_*.go registrations.
const (
	dashboardPath = "/" // GET /{$} exact root (see routes_dashboard.go)
	loginPath     = "/login"
	// loginRedirectAuth is the WithAuth redirect; must stay loginPath + ?error=1 (see paths_test).
	loginRedirectAuth = "/login?error=1"
	logoutPath        = "/logout"
	settingsPath      = "/settings"
	historyPath       = "/history" // must match [safepath.Default] (see paths_test).
	categoriesPath    = "/categories"
)
