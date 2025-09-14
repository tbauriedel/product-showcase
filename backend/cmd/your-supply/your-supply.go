package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tbauriedel/your-supply/internal/config"
	"github.com/tbauriedel/your-supply/internal/database"
	"github.com/tbauriedel/your-supply/internal/listener"
	"github.com/tbauriedel/your-supply/internal/version"
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

	// Setup context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var db database.Database

	switch cfg.Database.Type {
	case "mysql":
		db = database.NewMySQLDatabase(cfg, logger)
	default:
		logger.Error("Unknown database type provided. Check your configuration", "type", cfg.Database.Type)
		os.Exit(1)
	}

	// Check database connection
	if err = db.CheckStatus(ctx); err != nil {
		logger.Error("Could not connect to database", "error", err.Error())
	}

	// Start API listener with context
	if err := listener.New(cfg, logger, db).Run(ctx); err != nil {
		logger.Error("Listener finished with error", "error", err.Error())
	} else {
		logger.Info("Listener finished")
	}
}
