package main

import (
	"log"
	"net"
	"time"

	"context"
	"fmt"
	"os"

	"github.com/jansdhillon/task-manager/server/internal/server"
	"github.com/jansdhillon/task-manager/server/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
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

	store := &task.InMemoryTaskStore{
		Name:  "TaskManager",
		Tasks: make([]task.Task, 0, task.MAX_TASKS),
	}

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server.RegisterTaskService(s, store)
	reflection.Register(s)

	log.Printf("gRPC server listening on port 8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
