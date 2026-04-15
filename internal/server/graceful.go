package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ListenAndServeGracefully runs srv.ListenAndServe in the background, logs the listen address,
// blocks until SIGINT or SIGTERM, then shuts down with shutdownTimeout. If ListenAndServe fails
// with an error other than [http.ErrServerClosed], it logs and calls os.Exit(1).
func ListenAndServeGracefully(srv *http.Server, shutdownTimeout time.Duration) {
	slog.Info("listening", "addr", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server", "err", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(ctx)
	slog.Info("shutdown complete")
}
