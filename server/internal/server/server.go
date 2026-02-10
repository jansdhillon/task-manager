package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
	"github.com/jansdhillon/task-manager/server/internal/db"
)

type TaskServer struct {
	db *db.TaskDB
}

func NewTaskServer(database *db.TaskDB) *TaskServer {
	return &TaskServer{
		db: database,
	}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	task, err := s.db.CreateTask(ctx, req.Title, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &taskv1.CreateTaskResponse{
		Task: task,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := s.db.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &taskv1.GetTaskResponse{
		Task: task,
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := s.db.UpdateTask(ctx, id, req.Title, req.Description, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return &taskv1.UpdateTaskResponse{
		Task: task,
	}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*taskv1.DeleteTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	err = s.db.DeleteTask(ctx, id)
	if err != nil {
		return &taskv1.DeleteTaskResponse{
			Success: false,
		}, fmt.Errorf("failed to delete task: %w", err)
	}

	return &taskv1.DeleteTaskResponse{
		Success: true,
	}, nil
}
