package client

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/jansdhillon/task-manager/proto"
	"google.golang.org/grpc"
)

type mockTaskServiceClient struct {
	tasks map[string]*pb.Task
}

func (m *mockTaskServiceClient) CreateTask(ctx context.Context, req *pb.CreateTaskRequest, opts ...grpc.CallOption) (*pb.CreateTaskResponse, error) {
	task := &pb.Task{
		Id:    "test-id",
		Title: req.Title,
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	m.tasks[task.Id] = task
	return &pb.CreateTaskResponse{Task: task}, nil
}

func (m *mockTaskServiceClient) GetTask(ctx context.Context, req *pb.GetTaskRequest, opts ...grpc.CallOption) (*pb.GetTaskResponse, error) {
	task := m.tasks[req.Id]
	return &pb.GetTaskResponse{Task: task}, nil
}

func (m *mockTaskServiceClient) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest, opts ...grpc.CallOption) (*pb.UpdateTaskResponse, error) {
	task := m.tasks[req.Id]
	if task != nil {
		task.Title = req.Title
		if req.Description != nil {
			task.Description = req.Description
		}

		if req.Status != nil && *req.Status != task.Status {
			task.Status = *req.Status
		}
	}
	return &pb.UpdateTaskResponse{Task: task}, nil
}

func (m *mockTaskServiceClient) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest, opts ...grpc.CallOption) (*pb.DeleteTaskResponse, error) {
	delete(m.tasks, req.Id)
	return &pb.DeleteTaskResponse{Success: true}, nil
}

func TestTaskClient_CreateTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	title := "Test Task"
	description := "Test Description"

	task, err := client.CreateTask(ctx, title, &description)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if task.Title != title {
		t.Errorf("Expected title %s, got %s", title, task.Title)
	}
	if task.Description == nil || *task.Description != description {
		t.Errorf("Expected description %s, got %v", description, task.Description)
	}
}

func TestExecuteWithClient_CreateTask(t *testing.T) {
	originalNewTaskClient := NewTaskClient
	defer func() { NewTaskClient = originalNewTaskClient }()

	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	NewTaskClient = func(address string) (*TaskClient, error) {
		return &TaskClient{
			client: mockClient,
		}, nil
	}

	title := "Test Task"
	description := "Test Description"

	result, err := ExecuteWithClient("localhost:8080", func(c *TaskClient) (any, error) {
		return c.CreateTask(context.Background(), title, &description)
	})

	if err != nil {
		t.Fatalf("executeWithClient failed: %v", err)
	}

	task, ok := result.(*pb.Task)
	if !ok {
		t.Fatalf("Expected *pb.Task, got %T", result)
	}

	if task.Title != title {
		t.Errorf("Expected title %s, got %s", title, task.Title)
	}
	if task.Description == nil || *task.Description != description {
		t.Errorf("Expected description %s, got %v", description, task.Description)
	}
}

func TestExecuteWithClient_ConnectionError(t *testing.T) {
	originalNewTaskClient := NewTaskClient
	defer func() { NewTaskClient = originalNewTaskClient }()

	NewTaskClient = func(address string) (*TaskClient, error) {
		return nil, fmt.Errorf("connection failed")
	}

	result, err := ExecuteWithClient("localhost:8080", func(c *TaskClient) (any, error) {
		return c.CreateTask(context.Background(), "title", nil)
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if err.Error() != "connection failed" {
		t.Errorf("Expected 'connection failed', got %v", err)
	}
}

func TestExecuteWithClient_FunctionError(t *testing.T) {
	originalNewTaskClient := NewTaskClient
	defer func() { NewTaskClient = originalNewTaskClient }()

	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	NewTaskClient = func(address string) (*TaskClient, error) {
		return &TaskClient{
			client: mockClient,
		}, nil
	}

	result, err := ExecuteWithClient("localhost:8080", func(c *TaskClient) (any, error) {
		return nil, fmt.Errorf("function execution failed")
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
	if err.Error() != "function execution failed" {
		t.Errorf("Expected 'function execution failed', got %v", err)
	}
}
