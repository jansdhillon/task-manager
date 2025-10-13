package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
)

func GetDsn(ctx context.Context, sc secrets.SecretsClient) (string, error) {
	dsn, err := sc.GetLatestVersion(ctx, config.POSTGRES_DSN_ENV)
	dsn, err := sc.GetLatestVersion(ctx, config.POSTGRES_DSN_ENV)
	if err != nil {
		return "", err
	}
	dsnStr, ok := dsn.(string)
	if !ok {
		return "", fmt.Errorf("error converting DSN to a string: %s", dsnStr)
	}
	return dsn.(string), nil
}

func Connect() (*pgx.Conn, error) {
	projectId := os.Getenv(config.GCP_PROJECT_ID_ENV)
	if projectId == "" {
		return nil, errors.New("GCP project ID not found")
	}
	ctx := context.Background()
	sc, err := secrets.NewGcpSecretsClient(ctx, projectId)
	if err != nil {
		return nil, err
	}
	dsn, err := GetDsn(ctx, sc)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
