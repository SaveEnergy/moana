package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moana/internal/config"
	"moana/internal/db"
	"moana/internal/handlers"
	"moana/internal/store"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "user" {
		os.Exit(runUser(os.Args[2:]))
	}
	if len(os.Args) >= 2 && os.Args[1] == "serve" {
		runServe()
		return
	}
	if len(os.Args) >= 2 {
		fmt.Fprintf(os.Stderr, "usage: moana [serve | user <add|password> ...]\n")
		os.Exit(1)
	}
	runServe()
}

func runServe() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		slog.Error("database", "err", err)
		os.Exit(1)
	}
	defer database.Close()

	st := store.New(database)
	tmpl, err := handlers.ParseTemplates()
	if err != nil {
		slog.Error("templates", "err", err)
		os.Exit(1)
	}

	app := &handlers.App{
		Config:    cfg,
		Store:     st,
		Templates: tmpl,
	}

	h := handlers.NewRouter(app)
	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           h,
		ReadHeaderTimeout: cfg.RequestTimeout,
		ReadTimeout:       cfg.RequestTimeout * 2,
		WriteTimeout:      cfg.RequestTimeout * 2,
	}

	go func() {
		slog.Info("listening", "addr", cfg.Listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	slog.Info("shutdown complete")
}
