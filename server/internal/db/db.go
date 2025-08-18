package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jansdhillon/task-manager/server/internal/config"
)

func Connect() {
	conn, err := pgx.Connect(context.Background(), os.Getenv(config.DB_URL_ENV_VAR))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var id int
	var created_at time.Time
	err = conn.QueryRow(context.Background(), "select id, created_at from task").Scan(&id, &created_at)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(id, created_at)
}
