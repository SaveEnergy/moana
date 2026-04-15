package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"moana/internal/app"
	"moana/internal/config"
	"moana/internal/dbutil"
	"moana/internal/server"
)

func runServe() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	st, sqlDB, err := dbutil.OpenStore(cfg.DBPath)
	if err != nil {
		slog.Error("database", "err", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	handler, err := app.HTTPHandler(cfg, st)
	if err != nil {
		slog.Error("app", "err", err)
		os.Exit(1)
	}
	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           handler,
		ReadHeaderTimeout: cfg.RequestTimeout,
		ReadTimeout:       cfg.RequestTimeout * 2,
		WriteTimeout:      cfg.RequestTimeout * 2,
	}

	server.ListenAndServeGracefully(srv, 10*time.Second)
}
