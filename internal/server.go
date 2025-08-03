package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/jansdhillon/task-manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	store TaskStore
}

func NewTaskServer(store TaskStore) *TaskServer {
	return &TaskServer{
		store: store,
	}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	var description *string
	if req.Description != nil {
		description = req.Description
	}

	newTask := New(req.Title, description)
	createdTask, err := s.store.AddTask(newTask)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	pbTask := taskToProto(&createdTask)
	return &pb.CreateTaskResponse{
		Task: pbTask,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID: %w", err)
	}

	task, err := s.store.GetTask(id)
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

	updatedTask, err := s.store.UpdateTask(id, req.Title, description)
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

	err = s.store.DeleteTask(id)
	if err != nil {
		return &pb.DeleteTaskResponse{
			Success: false,
		}, fmt.Errorf("failed to delete task: %w", err)
	}

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

func taskToProto(t *Task) *pb.Task {
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

func RegisterTaskService(s *grpc.Server, store TaskStore) {
	pb.RegisterTaskServiceServer(s, NewTaskServer(store))
}
