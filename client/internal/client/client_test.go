package client

import (
	"context"
	"testing"

	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
)

type mockTaskServiceClient struct {
	createTaskFunc func(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error)
	getTaskFunc    func(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error)
	updateTaskFunc func(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error)
	deleteTaskFunc func(ctx context.Context, req *taskv1.DeleteTaskRequest) (*taskv1.DeleteTaskResponse, error)
}

func (m *mockTaskServiceClient) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.CreateTaskResponse, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, req)
	}
	return &taskv1.CreateTaskResponse{
		Task: &taskv1.Task{
			Id:    "test-id",
			Title: req.Title,
		},
	}, nil
}

func (m *mockTaskServiceClient) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.GetTaskResponse, error) {
	if m.getTaskFunc != nil {
		return m.getTaskFunc(ctx, req)
	}
	return &taskv1.GetTaskResponse{
		Task: &taskv1.Task{
			Id:    req.Id,
			Title: "Test Task",
		},
	}, nil
}

func (m *mockTaskServiceClient) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.UpdateTaskResponse, error) {
	if m.updateTaskFunc != nil {
		return m.updateTaskFunc(ctx, req)
	}
	return &taskv1.UpdateTaskResponse{
		Task: &taskv1.Task{
			Id:    req.Id,
			Title: req.Title,
		},
	}, nil
}

func (m *mockTaskServiceClient) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*taskv1.DeleteTaskResponse, error) {
	if m.deleteTaskFunc != nil {
		return m.deleteTaskFunc(ctx, req)
	}
	return &taskv1.DeleteTaskResponse{
		Success: true,
	}, nil
}

func TestTaskClient_CreateTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{}

	client := &TaskClient{
		client: mockClient,
	}

	desc := "test description"
	task, err := client.CreateTask(context.Background(), "Test Task", &desc)

	if err != nil {
		t.Errorf("CreateTask() unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("CreateTask() returned nil task")
	}

	if task.Title != "Test Task" {
		t.Errorf("CreateTask() task.Title = %v, want %v", task.Title, "Test Task")
	}
}

func TestTaskClient_GetTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{}

	client := &TaskClient{
		client: mockClient,
	}

	task, err := client.GetTask(context.Background(), "test-id")

	if err != nil {
		t.Errorf("GetTask() unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("GetTask() returned nil task")
	}

	if task.Id != "test-id" {
		t.Errorf("GetTask() task.Id = %v, want %v", task.Id, "test-id")
	}
}

func TestTaskClient_UpdateTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{}

	client := &TaskClient{
		client: mockClient,
	}

	task, err := client.UpdateTask(context.Background(), "test-id", "Updated Task", nil, nil)

	if err != nil {
		t.Errorf("UpdateTask() unexpected error: %v", err)
	}

	if task == nil {
		t.Fatal("UpdateTask() returned nil task")
	}

	if task.Title != "Updated Task" {
		t.Errorf("UpdateTask() task.Title = %v, want %v", task.Title, "Updated Task")
	}
}

func TestTaskClient_DeleteTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{}

	client := &TaskClient{
		client: mockClient,
	}

	success, err := client.DeleteTask(context.Background(), "test-id")

	if err != nil {
		t.Errorf("DeleteTask() unexpected error: %v", err)
	}

	if !success {
		t.Error("DeleteTask() success = false, want true")
	}
}
