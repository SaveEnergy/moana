package handlers

import "strconv"

// HTTP route paths for redirects and [net/http.ServeMux] registrations.
// Keep in sync with routes_*.go; exported for integration tests and external call sites.
const (
	// DashboardPath is the URL path for redirects to the app root.
	DashboardPath = "/"
	// DashboardRootPattern is the Go 1.22+ [http.ServeMux] pattern for an exact match on "/" only.
	// Use with [http.MethodGet]; bare "GET /" matches every path as a prefix — see routes_dashboard.go.
	DashboardRootPattern = "/{$}"
	LoginPath            = "/login"
	// ForgotPasswordPath and ResetPasswordPath back the password reset flow; see [registerAuthRoutes].
	ForgotPasswordPath = "/forgot-password"
	ResetPasswordPath  = "/reset-password"
	// PasswordResetQuerySent is set to 1 after a successful forgot-password POST (anti-enumeration copy).
	PasswordResetQuerySent = "sent"
	// LoginErrorQueryParam is the query key set when [WithAuth] redirects anonymous users ([LoginRedirectAuth]).
	LoginErrorQueryParam = "error"
	LoginRedirectAuth    = LoginPath + "?" + LoginErrorQueryParam + "=1" // WithAuth redirect; see paths_test.
	LogoutPath           = "/logout"
	SettingsPath         = "/settings"
	// SettingsErrorQueryParam and SettingsOKQueryParam are flash query keys for [App.Settings] (see settings_redirect).
	SettingsErrorQueryParam            = "err"
	SettingsOKQueryParam               = "ok"
	SettingsProfilePath                = SettingsPath + "/profile"
	// SettingsAvatarPath and SettingsAvatarRemovePath handle profile image upload/removal; see [registerSettingsRoutes] and [UserAvatarGet] for [AvatarsPathPattern].
	SettingsAvatarPath       = SettingsPath + "/avatar"
	SettingsAvatarRemovePath = SettingsAvatarPath + "/remove"
	// AvatarsPathPattern serves processed JPEGs for the signed-in user’s household; must match [UserAvatarGet].
	AvatarsPathPattern = "/avatars/{id}"
	SettingsHouseholdPath              = SettingsPath + "/household"
	SettingsHouseholdMembersPath       = SettingsHouseholdPath + "/members"
	SettingsHouseholdMembersRemovePath = SettingsHouseholdMembersPath + "/remove"
	HistoryPath                        = "/history" // must match [safepath.Default] (see paths_test).
	CategoriesPath                     = "/categories"
	CategoriesUpdatePath               = CategoriesPath + "/update"
	CategoriesDeletePath               = CategoriesPath + "/delete"
	TransactionsPath                   = "/transactions"
	// TransactionPathPattern and TransactionEditPathPattern are Go 1.22+ [http.ServeMux] patterns (wildcard segment {id}).
	TransactionPathPattern     = TransactionsPath + "/{id}"
	TransactionEditPathPattern = TransactionsPath + "/{id}/edit"
	// TransactionNextQueryParam is the form/query field for post-edit redirects ([safepath.Internal]); templates use the same name.
	TransactionNextQueryParam = "next"
	NotificationsPath         = "/notifications"
	// NotificationsMarkReadPath is POST target to mark one notification read.
	NotificationsMarkReadPath = NotificationsPath + "/read"
)

// TransactionURLPath returns the concrete URL path for one transaction (POST target, etc.).
// It matches [TransactionPathPattern] with this id substituted for {id}.
func TransactionURLPath(id int64) string {
	return TransactionsPath + "/" + strconv.FormatInt(id, 10)
}

// TransactionEditURLPath returns the concrete path for the edit form.
// It matches [TransactionEditPathPattern] with this id substituted for {id}.
func TransactionEditURLPath(id int64) string {
	return TransactionURLPath(id) + "/edit"
}
