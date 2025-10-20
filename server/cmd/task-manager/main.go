package main

import (
	"context"
	"log"
	"net"
	"os"

	"fmt"

	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/db"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
	"github.com/jansdhillon/task-manager/server/internal/server"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	gcpProjectId := os.Getenv("GCP_PROJECT_ID")

	if gcpProjectId == "" {
		log.Fatal("GCP project ID not found!")
	}
	ctx := context.Background()
	conn, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	defer conn.Close()

	sc, err := secrets.NewGcpSecretsClient(ctx, gcpProjectId)

	taskDB := db.NewTaskDB(conn, sc)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", config.SERVICE_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server.RegisterTaskService(s, taskDB)
	reflection.Register(s)

	log.Printf("gRPC server listening on port %s", config.SERVICE_PORT)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
