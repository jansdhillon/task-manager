package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
)

var (
	newSecretsClient = func(ctx context.Context, projectID string) (secrets.SecretsClient, error) {
		return secrets.NewGcpSecretsClient(ctx, projectID)
	}
	sqlOpenFunc = sql.Open
)

func GetDsn(ctx context.Context, sc secrets.SecretsClient) (string, error) {
	dsn, err := sc.GetLatestVersion(ctx, config.POSTGRES_DSN_ENV)
	if err != nil {
		return "", err
	}
	return dsn, nil
}

func Connect() (*sql.DB, error) {
	projectId := os.Getenv(config.GCP_PROJECT_ID_ENV)
	if projectId == "" {
		return nil, errors.New("GCP project ID not found")
	}
	ctx := context.Background()
	sc, err := newSecretsClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	dsn, err := GetDsn(ctx, sc)
	if err != nil {
		return nil, err
	}
	db, err := sqlOpenFunc("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return db, nil
}
