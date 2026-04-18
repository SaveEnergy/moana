package handlers

import "testing"

func TestFormFieldNameConstants_nonEmpty(t *testing.T) {
	t.Parallel()
	names := []string{
		LoginFieldEmail, LoginFieldPassword, LoginFieldRemember,
		SettingsFieldFirstName, SettingsFieldLastName,
		SettingsFieldCurrentPassword, SettingsFieldNewPassword, SettingsFieldNewPasswordConfirm,
		SettingsFieldHouseholdName,
	}
	for _, name := range names {
		if name == "" {
			t.Fatal("empty handler form field constant")
		}
	}
}
