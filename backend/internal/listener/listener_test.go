package listener

import (
	"context"
	"github.com/tbauriedel/your-supply/internal/config"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := config.Config{ListenAddr: ":1234"}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	l := New(cfg, logger)

	if l.config.ListenAddr != ":1234" {
		t.Fatalf("Actual: %s, Expected: %s", l.config.ListenAddr, ":1234")
	}

	if l.mux == nil {
		t.Fatal("Expected mux to be initialized, got nil")
	}
}

func TestRun(t *testing.T) {
	cfg := config.Config{ListenAddr: ":1234"}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	l := New(cfg, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		// Simulate SIGINT
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(os.Interrupt)
	}()

	err := l.Run(ctx)
	if err != nil {
		t.Errorf("Run returned error: %v", err)
	}
}

func TestHandleHealth(t *testing.T) {
	cfg := config.Config{ListenAddr: ":1234"}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	l := New(cfg, logger)

	req := httptest.NewRequest(http.MethodGet, "/-/health", nil)
	w := httptest.NewRecorder()

	l.handleHealth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Actual: %d, Expected: %d", resp.StatusCode, http.StatusOK)
	}
}
