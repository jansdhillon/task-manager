package main

import (
	"context"
	"log"
	"net"

	"fmt"

	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/db"
	"github.com/jansdhillon/task-manager/server/internal/server"
	"github.com/jansdhillon/task-manager/server/internal/task"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	store := &task.InMemoryTaskStore{
		Name:  "TaskManager",
		Tasks: make([]*task.Task, 0, task.MAX_TASKS),
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	defer conn.Close(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", config.SERVICE_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server.RegisterTaskService(s, store)
	reflection.Register(s)

	log.Printf("gRPC server listening on port %s", config.SERVICE_PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
