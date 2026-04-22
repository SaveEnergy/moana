package handlers

// HTML form control names for login and settings POST handlers (keep in sync with assets templates).
// Household member add uses [LoginFieldEmail] and [LoginFieldPassword] (same names as the sign-in form).
const (
	LoginFieldEmail    = "email"
	LoginFieldPassword = "password"
	LoginFieldRemember = "remember"
)

// Password reset: forgot-password and reset-password forms ([forgot_password.html], [reset_password.html]).
const (
	// ResetFieldToken is the query key and form field for the one-time token.
	ResetFieldToken = "token"
	// ResetFieldPassword reuses the same control name as [LoginFieldPassword] ("password").
	ResetFieldPasswordConfirm = "password_confirm"
)

const (
	SettingsFieldFirstName          = "first_name"
	SettingsFieldLastName           = "last_name"
	SettingsFieldCurrentPassword    = "current_password"
	SettingsFieldNewPassword        = "new_password"
	SettingsFieldNewPasswordConfirm = "new_password_confirm"
	SettingsFieldHouseholdName      = "household_name"
	// SettingsFieldMemberUserID is the hidden field for member remove / leave (settings.html).
	SettingsFieldMemberUserID = "user_id"
	// SettingsFieldAvatar is the file field name for POST /settings/avatar.
	SettingsFieldAvatar = "avatar"
	// NotificationFieldID is the notification id for POST /notifications/read (notifications.html).
	NotificationFieldID = "id"
)
