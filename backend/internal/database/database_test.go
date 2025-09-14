package database

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestCheckStatus(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	defer db.Close()

	// Mock database connection and expect ping to be called
	sqlDB := sqlx.NewDb(db, "sqlmock")
	mock.ExpectPing().WillReturnError(nil)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	s := &SQLDatabase{
		db:     sqlDB,
		logger: logger,
	}

	err := s.CheckStatus(context.TODO())
	if err != nil {
		t.Fatalf("error checking status: %v", err)
	}
}
