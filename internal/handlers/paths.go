package handlers

// HTTP route paths for redirects and [net/http.ServeMux] registrations.
// Keep in sync with routes_*.go; exported for integration tests and external call sites.
const (
	DashboardPath = "/" // GET /{$} exact root (see routes_dashboard.go)
	LoginPath     = "/login"
	// LoginRedirectAuth is WithAuth's redirect; must stay LoginPath + ?error=1 (see paths_test).
	LoginRedirectAuth = "/login?error=1"
	LogoutPath        = "/logout"
	SettingsPath      = "/settings"
	// Settings POST targets (prefix SettingsPath).
	SettingsProfilePath                = "/settings/profile"
	SettingsHouseholdPath              = "/settings/household"
	SettingsHouseholdMembersPath       = "/settings/household/members"
	SettingsHouseholdMembersRemovePath = "/settings/household/members/remove"
	HistoryPath                        = "/history" // must match [safepath.Default] (see paths_test).
	CategoriesPath                     = "/categories"
	CategoriesUpdatePath               = "/categories/update"
	CategoriesDeletePath               = "/categories/delete"
	TransactionsPath                   = "/transactions"
	NotificationsPath                  = "/notifications"
)
