package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type DSN struct {
	Database string
	Username string
	Password string
	Host     string
	Port     string
	Schema   string
}

type DSNOpt func(*DSN)

func NewDSN(opts ...DSNOpt) *DSN {
	dsn := &DSN{
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
		Username: "postgres",
		Password: "password",
		Schema:   "public",
	}
	for _, o := range opts {
		o(dsn)
	}

	return dsn
}

func WithHost(host string) DSNOpt {
	return func(d *DSN) {
		if host != "" {
			d.Host = host
		}
	}
}

func WithDatabase(database string) DSNOpt {
	return func(d *DSN) {
		if database != "" {
			d.Database = database
		}
	}
}

func WithUsername(username string) DSNOpt {
	return func(d *DSN) {
		if username != "" {
			d.Username = username
		}
	}
}

func WithPassword(password string) DSNOpt {
	return func(d *DSN) {
		if password != "" {
			d.Password = password
		}
	}
}

func WithPort(port string) DSNOpt {
	return func(d *DSN) {
		if port != "" {
			d.Port = port
		}
	}
}

func (d *DSN) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", d.Username, d.Password, d.Host, d.Port, d.Database)
}

func Connect(ctx context.Context, opts ...DSNOpt) (*pgx.Conn, error) {
	defaultOpts := []DSNOpt{
		WithHost(os.Getenv("TASK_MANAGER_DB_HOST")),
		WithPort(os.Getenv("TASK_MANAGER_DB_PORT")),
		WithDatabase(os.Getenv("TASK_MANAGER_DB_NAME")),
		WithUsername(os.Getenv("TASK_MANAGER_DB_USER")),
		WithPassword(os.Getenv("TASK_MANAGER_DB_PASSWORD")),
	}
	allOpts := append(defaultOpts, opts...)

	dsn := NewDSN(allOpts...)
	conn, err := pgx.Connect(ctx, dsn.String())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	defer conn.Close(ctx)
	return conn, nil
}
