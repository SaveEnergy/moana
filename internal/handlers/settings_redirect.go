package handlers

import (
	"net/http"
	"net/url"
)

func redirectSettingsErr(w http.ResponseWriter, r *http.Request, msg string) {
	http.Redirect(w, r, SettingsPath+"?"+SettingsErrorQueryParam+"="+url.QueryEscape(msg), http.StatusSeeOther)
}

func redirectSettingsOK(w http.ResponseWriter, r *http.Request, okKey string) {
	http.Redirect(w, r, SettingsPath+"?"+SettingsOKQueryParam+"="+url.QueryEscape(okKey), http.StatusSeeOther)
}
