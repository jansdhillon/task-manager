package main

import (
	"log"
	"net"

	. "github.com/jansdhillon/task-manager/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	store := &InMemoryTaskStore{
		Name:  "TaskManager",
		Tasks: make([]Task, 0, MAX_TASKS),
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterTaskService(s, store)
	reflection.Register(s)

	log.Printf("gRPC server listening on :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
