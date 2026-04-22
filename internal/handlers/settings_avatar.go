package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"moana/internal/avatar"
	"moana/internal/httperr"
	"moana/internal/store"
)

// SettingsAvatarUpload handles POST /settings/avatar (multipart; field "avatar").
func (a *App) SettingsAvatarUpload(w http.ResponseWriter, r *http.Request, u *store.User) {
	if err := r.ParseMultipartForm(4 * 1 << 20); err != nil {
		redirectSettingsErr(w, r, "Could not read upload.")
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()
	file, fileHeader, err := r.FormFile(SettingsFieldAvatar)
	if err != nil {
		redirectSettingsErr(w, r, "Choose an image file.")
		return
	}
	defer func() { _ = file.Close() }()
	if a.AvatarDir == "" {
		httperr.Internal(w, r, errors.New("avatar directory not configured"))
		return
	}
	raw, err := readAvatarFile(file, fileHeader, avatar.MaxUploadBytes)
	if err != nil {
		redirectSettingsErr(w, r, "Image file is too large (max 1 MB).")
		return
	}
	if len(raw) == 0 {
		redirectSettingsErr(w, r, "Choose a non-empty image file.")
		return
	}
	jpegBytes, err := avatar.ToJPEG(raw)
	if err != nil {
		switch {
		case errors.Is(err, avatar.ErrEmpty), errors.Is(err, avatar.ErrInvalid):
			redirectSettingsErr(w, r, "That file is not a valid image (PNG, JPEG, GIF, or WebP).")
		case errors.Is(err, avatar.ErrTooLarge), errors.Is(err, avatar.ErrTooBig):
			redirectSettingsErr(w, r, "Image is too large. Try a smaller or shorter file.")
		default:
			httperr.Internal(w, r, err)
		}
		return
	}
	dest := userAvatarFilePath(a.AvatarDir, u.ID)
	tmp, err := os.CreateTemp(a.AvatarDir, "avatar-incomplete-*")
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	tmpPath := tmp.Name()
	_, werr := tmp.Write(jpegBytes)
	if werr == nil {
		werr = tmp.Close()
	} else {
		_ = tmp.Close()
	}
	if werr != nil {
		_ = os.Remove(tmpPath)
		httperr.Internal(w, r, werr)
		return
	}
	if err := os.Rename(tmpPath, dest); err != nil {
		_ = os.Remove(tmpPath)
		httperr.Internal(w, r, err)
		return
	}
	if err := a.Store.IncrementUserAvatarRev(r.Context(), u.ID); err != nil {
		_ = os.Remove(dest)
		httperr.Internal(w, r, err)
		return
	}
	redirectSettingsOK(w, r, "avatar")
}

// SettingsAvatarRemove handles POST /settings/avatar/remove (own photo only).
func (a *App) SettingsAvatarRemove(w http.ResponseWriter, r *http.Request, u *store.User) {
	if !requireParseFormSettings(w, r) {
		return
	}
	if a.AvatarDir == "" {
		httperr.Internal(w, r, errors.New("avatar directory not configured"))
		return
	}
	_ = os.Remove(userAvatarFilePath(a.AvatarDir, u.ID))
	if err := a.Store.ClearUserAvatar(r.Context(), u.ID); err != nil {
		httperr.Internal(w, r, err)
		return
	}
	redirectSettingsOK(w, r, "avatar-removed")
}

// UserAvatarGet serves a household-scoped profile JPEG (GET /avatars/{id}).
func (a *App) UserAvatarGet(w http.ResponseWriter, r *http.Request, u *store.User) {
	if a.AvatarDir == "" {
		httperr.Internal(w, r, errors.New("avatar directory not configured"))
		return
	}
	uid, ok := pathPositiveInt64(r, "id")
	if !ok {
		http.NotFound(w, r)
		return
	}
	target, err := a.Store.GetUserByID(r.Context(), uid)
	if err != nil {
		httperr.Internal(w, r, err)
		return
	}
	if target == nil || target.AvatarRev == 0 {
		http.NotFound(w, r)
		return
	}
	if target.HouseholdID != u.HouseholdID {
		http.NotFound(w, r)
		return
	}
	p := userAvatarFilePath(a.AvatarDir, uid)
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		httperr.Internal(w, r, err)
		return
	}
	defer func() { _ = f.Close() }()
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "private, max-age=3600")
	_, _ = io.Copy(w, f)
}

func readAvatarFile(r io.Reader, h *multipart.FileHeader, max int64) ([]byte, error) {
	lim := max + 1
	if h != nil && h.Size > 0 && h.Size < lim {
		lim = h.Size
	}
	b, err := io.ReadAll(io.LimitReader(r, lim))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > max {
		return nil, errors.New("file too long")
	}
	return b, nil
}

func userAvatarFilePath(avatarDir string, userID int64) string {
	return filepath.Join(avatarDir, fmt.Sprintf("%d.jpg", userID))
}
