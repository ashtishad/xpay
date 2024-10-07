package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// NewConnection creates a new db connection using pgx driver, return a standard *sql.DB handle.
// It applies connection pool settings and verifies the connection with a ping.
func NewConnection(ctx context.Context, cfg common.DBConfig, l *slog.Logger) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid connection string: %w", err)
	}

	db := stdlib.OpenDB(*connConfig)

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}

	l.Info("successfully connected to postgres", "dsn", connConfig.ConnString())

	return db, nil
}
