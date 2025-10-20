package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
)

const defaultDriver = "pgx"

func GetDsn(ctx context.Context, sc secrets.SecretsClient) (string, error) {
	dsn, err := sc.GetLatestVersion(ctx, config.POSTGRES_DSN_ENV)
	if err != nil {
		return "", err
	}
	return dsn, nil
}

// Connect to the DSN with the default driver
func Connect(ctx context.Context, dsn string) (*sql.DB, error) {
	return connect(ctx, defaultDriver, dsn)
}

func connect(ctx context.Context, driver, dsn string) (*sql.DB, error) {
	if driver == "" {
		driver = defaultDriver
	}

	if driver == defaultDriver {
		if cfg, err := pgx.ParseConfig(dsn); err == nil {
			cfg.StatementCacheCapacity = 0
			cfg.DescriptionCacheCapacity = 0
			cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
			dsn = stdlib.RegisterConnConfig(cfg)
		}
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}
