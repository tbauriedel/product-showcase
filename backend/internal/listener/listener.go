package listener

import (
	"context"
	"encoding/json"
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
	"github.com/tbauriedel/product-showcase/internal/model"
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
	mux.HandleFunc("GET /health", l.handleHealth)

	mux.HandleFunc("POST /v1/product", l.handleInsertProduct)

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

// handleInsertProduct handles the insertion of a new product into the database
func (l *Listener) handleInsertProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		http.Error(w, fmt.Sprint("invalid json provided. ", err), http.StatusBadRequest)
		return
	}

	// TODO validate provided fields

	// Save into database
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = l.DB.InsertProduct(ctx, &product)
	if err != nil {
		http.Error(w, fmt.Sprint("could not insert product into database. ", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
