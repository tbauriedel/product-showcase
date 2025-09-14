package database

import (
	"fmt"
	"log/slog"

	"github.com/tbauriedel/your-supply/internal/config"
)

type MySQLDatabase struct {
	SQLDatabase
}

// NewMySQLDatabase creates a new MySQLDatabase instance and returns it
func NewMySQLDatabase(cfg config.Config, logger *slog.Logger) *MySQLDatabase {
	database := &MySQLDatabase{
		SQLDatabase: NewSQLDatabase(cfg, logger, "mysql", getDnsMySQL(cfg)),
	}

	return database
}

// getDnsMySQL returns the DSN string for MySQL connection based on the given config.Config
func getDnsMySQL(cfg config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Address,
		cfg.Database.Port,
		cfg.Database.DbName)
}
