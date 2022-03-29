package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hedwig100/bookmark/backend/server"
	"github.com/hedwig100/bookmark/backend/slog"
)

func main() {
	server := server.Server()

	// For graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sig

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		slog.Info("Server shutdown...")
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			slog.Errf("HTTP server Shutdown: %v", err)
		}
		slog.Info("Sucessful in graceful shutdown!")
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Fatalf("server listening failed %v", err)
	}

	<-idleConnsClosed
}
