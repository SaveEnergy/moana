package handlers

import (
	"html/template"
	"moana/internal/config"
	"moana/internal/store"
)

// App wires HTTP handlers to shared dependencies.
type App struct {
	Config    *config.Config
	Store     *store.Store
	Templates *template.Template
}
