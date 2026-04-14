package household

import (
	"context"
	"errors"

	"moana/internal/store"
)

// ErrHouseholdMissing is returned when the user's household row is absent.
var ErrHouseholdMissing = errors.New("household not found")

// SettingsPageData is the template payload for the settings screen (profile + household).
type SettingsPageData struct {
	Error            string
	Success          string
	User             *store.User
	Household        *store.Household
	Members          []store.HouseholdMember
	MemberCount      int64
	CanManageMembers bool
}

// LoadSettingsPage loads household rows and permission flags for the settings UI.
func LoadSettingsPage(ctx context.Context, st *store.Store, u *store.User, errMsg, successMsg string) (SettingsPageData, error) {
	hh, err := st.GetHousehold(ctx, u.HouseholdID)
	if err != nil {
		return SettingsPageData{}, err
	}
	if hh == nil {
		return SettingsPageData{}, ErrHouseholdMissing
	}
	members, err := st.ListHouseholdMembers(ctx, u.HouseholdID)
	if err != nil {
		return SettingsPageData{}, err
	}
	n, err := st.CountHouseholdMembers(ctx, u.HouseholdID)
	if err != nil {
		return SettingsPageData{}, err
	}
	return SettingsPageData{
		Error:            errMsg,
		Success:          successMsg,
		User:             u,
		Household:        hh,
		Members:          members,
		MemberCount:      n,
		CanManageMembers: CanManage(u),
	}, nil
}
