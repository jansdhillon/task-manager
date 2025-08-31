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

func getDsn(ctx context.Context, sc secrets.SecretsClient) (string, error) {
	password, err := sc.GetLatestVersion(ctx, "POSTGRES_PASSWORD")
	if err != nil {
		return "", err
	}
	host := "aws-0-ca-central-1.pooler.supabase.com"
	port := 6543
	database := "postgres"
	user := "postgres.iarmdhqsjqxbirfugjjc"

	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s", database, user, password, host, port, database)

	return dsn, nil
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
	dsn, err := getDsn(ctx, sc)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
