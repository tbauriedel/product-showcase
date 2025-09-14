package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/tbauriedel/product-showcase/internal/config"
	"github.com/tbauriedel/product-showcase/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Database is the interface that wraps the db functions
type Database interface {
	CheckStatus(ctx context.Context) error
	InsertProduct(ctx context.Context, product *model.Product) error
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

	s.logger.Debug("Pinged database. Connection OK", "type", s.config.Database.Type)

	return nil
}

// InsertProduct inserts a new product into the database
//
// Returns error if insertion fails
func (s *SQLDatabase) InsertProduct(ctx context.Context, product *model.Product) error {
	query := s.db.Rebind("INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)")

	result, err := s.db.ExecContext(ctx,
		query,
		product.Name, product.Description, product.Price, product.Stock, time.Now(), time.Now(),
	)

	if err != nil {
		return fmt.Errorf("error inserting product: %w", err)
	}

	id, _ := result.LastInsertId()
	s.logger.Debug("Inserted product into database", "id", id)

	return nil
}
