package handlers

// HTTP route paths for redirects and [net/http.ServeMux] registrations.
// Keep in sync with routes_*.go; exported for integration tests and external call sites.
const (
	// DashboardPath is the URL path for redirects to the app root.
	DashboardPath = "/"
	// DashboardRootPattern is the Go 1.22+ [http.ServeMux] pattern for an exact match on "/" only.
	// Use with [http.MethodGet]; bare "GET /" matches every path as a prefix — see routes_dashboard.go.
	DashboardRootPattern               = "/{$}"
	LoginPath                          = "/login"
	LoginRedirectAuth                  = LoginPath + "?error=1" // WithAuth redirect; see paths_test.
	LogoutPath                         = "/logout"
	SettingsPath                       = "/settings"
	SettingsProfilePath                = SettingsPath + "/profile"
	SettingsHouseholdPath              = SettingsPath + "/household"
	SettingsHouseholdMembersPath       = SettingsHouseholdPath + "/members"
	SettingsHouseholdMembersRemovePath = SettingsHouseholdMembersPath + "/remove"
	HistoryPath                        = "/history" // must match [safepath.Default] (see paths_test).
	CategoriesPath                     = "/categories"
	CategoriesUpdatePath               = CategoriesPath + "/update"
	CategoriesDeletePath               = CategoriesPath + "/delete"
	TransactionsPath                   = "/transactions"
	NotificationsPath                  = "/notifications"
)
