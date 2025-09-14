package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tbauriedel/product-showcase/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Database is the interface that wraps the db functions
type Database interface {
	CheckStatus(ctx context.Context) error
}

// SQLDatabase represents a SQL database connection
type SQLDatabase struct {
	db         *sqlx.DB
	config     config.Config
	logger     *slog.Logger
	connectErr error
}

// NewSQLDatabase creates a new SQLDatabase instance and returns it
//
// Errors during connection are stored in the struct and can be checked later
func NewSQLDatabase(cfg config.Config, logger *slog.Logger, driver string, dsn string) SQLDatabase {
	db, err := sqlx.Connect(driver, dsn)

	database := SQLDatabase{
		db:         db,
		config:     cfg,
		logger:     logger,
		connectErr: err,
	}

	return database
}

// CheckStatus checks the status of the database connection
//
// Pings database and returns error if ping fails
func (s *SQLDatabase) CheckStatus(ctx context.Context) error {
	if s.connectErr != nil {
		return fmt.Errorf("error opening database connection: %w", s.connectErr)
	}

	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.logger.Debug("Pinged database. Connection OK")

	return nil
}
