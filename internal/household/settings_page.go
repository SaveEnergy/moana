package household

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

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
// GetHousehold and ListHouseholdMembers run concurrently (independent reads on the same DB).
func LoadSettingsPage(ctx context.Context, st *store.Store, u *store.User, errMsg, successMsg string) (SettingsPageData, error) {
	var (
		hh      *store.Household
		members []store.HouseholdMember
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		hh, err = st.GetHousehold(gctx, u.HouseholdID)
		return err
	})
	g.Go(func() error {
		var err error
		members, err = st.ListHouseholdMembers(gctx, u.HouseholdID)
		return err
	})
	if err := g.Wait(); err != nil {
		return SettingsPageData{}, err
	}
	if hh == nil {
		return SettingsPageData{}, ErrHouseholdMissing
	}
	// MemberCount matches the list length (same as COUNT(*) for this household) without a second query.
	n := int64(len(members))
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
