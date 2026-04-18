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

func TestFormFieldNameConstants_profilePasswordKeysDistinct(t *testing.T) {
	t.Parallel()
	a, b, c := SettingsFieldCurrentPassword, SettingsFieldNewPassword, SettingsFieldNewPasswordConfirm
	if a == b || a == c || b == c {
		t.Fatalf("profile password keys must be distinct, got %q %q %q", a, b, c)
	}
}
