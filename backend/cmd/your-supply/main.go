package main

import (
	"context"
	"flag"
	"github.com/tbauriedel/your-supply/internal/version"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tbauriedel/your-supply/internal/config"
	"github.com/tbauriedel/your-supply/internal/listener"
)

func main() {
	var (
		flagConfigFile string
		flagVersion    bool
	)

	flag.StringVar(&flagConfigFile, "config", "", "path to config file")
	flag.BoolVar(&flagVersion, "version", false, "print version and exit")

	flag.Parse()

	if flagVersion {
		println("Your-Supply version:", version.Version)
		os.Exit(0)
	}

	// Get config from file
	cfg, err := config.NewFromFile(flagConfigFile)
	if err != nil {
		slog.Error("Error loading config", "error", err.Error())
		os.Exit(1)
	}

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.GetLogLevel()}))

	// Setup context with signal handling. Needed for listener
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start API listener with context
	if err := listener.New(cfg, logger).Run(ctx); err != nil {
		logger.Error("Listener finished with error", "error", err.Error())
	} else {
		logger.Info("Listener finished")
	}
}
