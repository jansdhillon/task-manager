package server

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	pb "github.com/jansdhillon/task-manager/proto"
	"github.com/jansdhillon/task-manager/server/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	store task.TaskStore
}

func NewTaskServer(store task.TaskStore) *TaskServer {
	return &TaskServer{
		store: store,
	}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	createdTask, err := s.store.AddTask(ctx, req.Title, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	pbTask := taskToProto(createdTask)
	log.Printf("task created: %v", pbTask)
	return &pb.CreateTaskResponse{
		Task: pbTask,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := s.store.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	pbTask := taskToProto(task)
	return &pb.GetTaskResponse{
		Task: pbTask,
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	updatedTask, err := s.store.UpdateTask(ctx, id, req.Title, description, task.InProgress)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	pbTask := taskToProto(updatedTask)
	return &pb.UpdateTaskResponse{
		Task: pbTask,
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	err = s.store.DeleteTask(ctx, id)
	if err != nil {
		return &pb.DeleteTaskResponse{
			Success: false,
		}, fmt.Errorf("failed to delete task: %w", err)
	}

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

func taskToProto(t *task.Task) *pb.Task {
	pbTask := &pb.Task{
		Id:            t.ID.String(),
		Title:         t.Title,
		CreatedAt:     timestamppb.New(t.CreatedAt),
		LastUpdatedAt: timestamppb.New(t.LastUpdatedAt),
		Deleted:       t.Deleted,
	}

	if t.Description != nil {
		pbTask.Description = t.Description
	}

	return pbTask
}

func RegisterTaskService(s *grpc.Server, store task.TaskStore) {
	pb.RegisterTaskServiceServer(s, NewTaskServer(store))
}
