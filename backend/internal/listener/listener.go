package listener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tbauriedel/product-showcase/internal/config"
	"github.com/tbauriedel/product-showcase/internal/database"
	"github.com/tbauriedel/product-showcase/internal/version"
)

type Listener struct {
	config config.Config
	mux    http.Handler
	logger *slog.Logger
	DB     database.Database
}

// New creates a new Listener instance with the provided configuration and logger.
func New(config config.Config, logger *slog.Logger, db database.Database) *Listener {
	l := &Listener{
		config: config,
		logger: logger,
		DB:     db,
	}

	mux := http.NewServeMux()

	// Add routes here
	mux.HandleFunc("GET /-/health", l.handleHealth)

	l.mux = mux

	return l
}

// Run starts the HTTP server and listens for shutdown signals.
func (l *Listener) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:        l.config.ListenAddr,
		Handler:     l.mux,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	go func() {
		l.logger.Info("Started listener", slog.String("addr", l.config.ListenAddr))

		http.ListenAndServe(l.config.ListenAddr, l.mux)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.logger.Error("HTTP server error", "error", err.Error())
		}

		l.logger.Info("Received shutdown signal, shutting down listener")
	}()

	// Shutdown with timeout on SIGINT or SIGTERM. Timeout to shut down the application gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 3*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error", "error", err.Error())
	}

	l.logger.Info("Completed shutdown")

	return nil
}

// handleHealth responds with a simple health check message.
func (l *Listener) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	dbState := "ok"
	err := l.DB.CheckStatus(ctx)

	if err != nil {
		dbState = err.Error()
	}

	msg := fmt.Sprintf(`{"database": "%s", "version": "%s"}`, dbState, version.Version)

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, msg)
}
