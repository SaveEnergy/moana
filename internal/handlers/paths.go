package handlers

// HTTP route paths for redirects and [net/http.ServeMux] registrations.
// Keep in sync with routes_*.go; exported for integration tests and external call sites.
const (
	DashboardPath                      = "/" // GET /{$} exact root (see routes_dashboard.go)
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
