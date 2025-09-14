package database

import (
	"log/slog"

	"github.com/tbauriedel/product-showcase/internal/config"
)

type SQLiteDatabase struct {
	SQLDatabase
}

// NewSQLiteDatabase creates a new SQLiteDatabase instance and returns it
//
// We use the Address field of the config.Database struct to get the path to the SQLite database file
func NewSQLiteDatabase(cfg config.Config, logger *slog.Logger) *SQLiteDatabase {
	database := &SQLiteDatabase{
		SQLDatabase: NewSQLDatabase(cfg, logger, "sqlite3", cfg.Database.Address),
	}

	return database
}
