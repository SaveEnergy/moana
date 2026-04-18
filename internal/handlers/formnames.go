package handlers

// HTML form control names for login and settings POST handlers (keep in sync with assets templates).
// Household member add uses [LoginFieldEmail] and [LoginFieldPassword] (same names as the sign-in form).
const (
	LoginFieldEmail    = "email"
	LoginFieldPassword = "password"
	LoginFieldRemember = "remember"
)

const (
	SettingsFieldFirstName          = "first_name"
	SettingsFieldLastName           = "last_name"
	SettingsFieldCurrentPassword    = "current_password"
	SettingsFieldNewPassword        = "new_password"
	SettingsFieldNewPasswordConfirm = "new_password_confirm"
	SettingsFieldHouseholdName      = "household_name"
)
