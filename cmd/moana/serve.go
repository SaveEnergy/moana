package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"moana/internal/app"
	"moana/internal/config"
	"moana/internal/dbutil"
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
