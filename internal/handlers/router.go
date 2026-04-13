package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"moana/internal/store"
)

// RouterOptions configures HTTP routing (e.g. for tests).
type RouterOptions struct {
	DisableRequestLogging bool
}

// NewRouter registers all application routes on mux.
func NewRouter(app *App) http.Handler {
	return NewRouterWithRouterOptions(nil, app)
}

// NewRouterWithRouterOptions registers routes with optional logging disabled (integration tests).
func NewRouterWithRouterOptions(opts *RouterOptions, app *App) http.Handler {
	mux := http.NewServeMux()

	staticFS, err := StaticFS()
	if err != nil {
		panic(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /login", app.LoginPage)
	mux.HandleFunc("POST /login", app.LoginSubmit)
	mux.HandleFunc("POST /logout", app.Logout)

	mux.Handle("GET /", app.WithAuth(app.Home))
	mux.Handle("GET /transactions", app.WithAuth(app.Transactions))
	mux.Handle("POST /transactions", app.WithAuth(app.TransactionCreate))
	mux.Handle("GET /transactions/{id}/edit", app.WithAuth(app.TransactionEdit))
	mux.Handle("POST /transactions/{id}", app.WithAuth(app.TransactionUpdate))
	mux.Handle("GET /history", app.WithAuth(app.History))
	mux.Handle("GET /categories", app.WithAuth(app.Categories))
	mux.Handle("POST /categories", app.WithAuth(app.CategoryCreate))
	mux.Handle("POST /categories/delete", app.WithAuth(app.CategoryDelete))

	mux.Handle("GET /admin/users", app.WithAuthAdmin(app.AdminUsers))
	mux.Handle("POST /admin/users", app.WithAuthAdmin(app.AdminUserCreate))
	mux.Handle("POST /admin/users/password", app.WithAuthAdmin(app.AdminUserPassword))

	var h http.Handler = mux
	if opts == nil || !opts.DisableRequestLogging {
		h = loggingMiddleware(mux)
	}
	return h
}

// WithAuth requires a valid session and loads the current user.
func (a *App) WithAuth(next func(http.ResponseWriter, *http.Request, *store.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.currentUser(r)
		if err != nil || u == nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		next(w, r, u)
	})
}

// WithAuthAdmin requires a signed-in user with role admin.
func (a *App) WithAuthAdmin(next func(http.ResponseWriter, *http.Request, *store.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.currentUser(r)
		if err != nil || u == nil {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
			return
		}
		if u.Role != "admin" {
			http.Error(w, "Administrator access required.", http.StatusForbidden)
			return
		}
		next(w, r, u)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lw.status,
			"dur_ms", time.Since(start).Milliseconds(),
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
