package handlers

import (
	"moana/internal/config"
	"moana/internal/render"
	"moana/internal/store"
)

// App wires HTTP handlers to shared dependencies. Prefer moana/internal/app.New for a fully wired instance.
type App struct {
	Config *config.Config
	Store  *store.Store
	Render *render.Engine
}
