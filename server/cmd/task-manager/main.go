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
	projectId := os.Getenv("GCP_PROJECT_ID")

	if projectId == "" {
		log.Fatal("GCP project ID not found!")
	}
	ctx := context.Background()
	sc, err := secrets.NewGcpSecretsClient(ctx, projectId)
	if err != nil {
		log.Fatalf("error creating secrets client: %v", err)
	}
	dsn, err := db.GetDsn(ctx, sc)
	if err != nil {
		log.Fatalf("error getting dsn: %v", err)
	}
	conn, err := db.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	defer conn.Close()

	taskDB := db.NewTaskDB(conn)

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
