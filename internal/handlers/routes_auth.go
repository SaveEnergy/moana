package handlers

import "net/http"

func registerAuthRoutes(mux *http.ServeMux, app *App) {
	mux.HandleFunc(http.MethodGet+" "+LoginPath, app.LoginPage)
	mux.HandleFunc(http.MethodPost+" "+LoginPath, app.LoginSubmit)
	mux.HandleFunc(http.MethodGet+" "+ForgotPasswordPath, app.ForgotPasswordPage)
	mux.HandleFunc(http.MethodPost+" "+ForgotPasswordPath, app.ForgotPasswordSubmit)
	mux.HandleFunc(http.MethodGet+" "+ResetPasswordPath, app.ResetPasswordPage)
	mux.HandleFunc(http.MethodPost+" "+ResetPasswordPath, app.ResetPasswordSubmit)
	mux.HandleFunc(http.MethodPost+" "+LogoutPath, app.Logout)
}
