package client

import (
	"context"
	"fmt"
	"strings"
	"testing"

	pb "github.com/jansdhillon/task-manager/proto"
	"google.golang.org/grpc"
)

type mockTaskServiceClient struct {
	tasks     map[string]*pb.Task
	shouldErr bool
}

func (m *mockTaskServiceClient) CreateTask(ctx context.Context, req *pb.CreateTaskRequest, opts ...grpc.CallOption) (*pb.CreateTaskResponse, error) {
	if m.shouldErr {
		return nil, fmt.Errorf("mock error")
	}
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
	if m.shouldErr {
		return nil, fmt.Errorf("mock error")
	}
	task := m.tasks[req.Id]
	return &pb.GetTaskResponse{Task: task}, nil
}

func (m *mockTaskServiceClient) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest, opts ...grpc.CallOption) (*pb.UpdateTaskResponse, error) {
	if m.shouldErr {
		return nil, fmt.Errorf("mock error")
	}
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
	if m.shouldErr {
		return nil, fmt.Errorf("mock error")
	}
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

func TestUpdateWithClient(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	
	oldDesc := "old description"
	newDesc := "new description"
	
	tt := []struct {
		testName      string
		title         struct{ oldTitle, newTitle string }
		description   struct{ oldDescription, newDescription *string }
	}{
		{
			testName: "update title",
			title: struct {
				oldTitle string
				newTitle string
			}{oldTitle: "old title", newTitle: "new title"},
		},
		{
			testName: "update description",
			title: struct {
				oldTitle string
				newTitle string
			}{oldTitle: "test title", newTitle: "test title"},
			description: struct {
				oldDescription  *string
				newDescription *string
			}{oldDescription: &oldDesc, newDescription: &newDesc},
		},
		{
			testName: "update both title and description",
			title: struct {
				oldTitle string
				newTitle string
			}{oldTitle: "old title", newTitle: "updated title"},
			description: struct {
				oldDescription  *string
				newDescription *string
			}{oldDescription: &oldDesc, newDescription: &newDesc},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			originalTask, err := client.CreateTask(ctx, tc.title.oldTitle, tc.description.oldDescription)
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			updatedTask, err := client.UpdateTask(ctx, originalTask.Id, tc.title.newTitle, tc.description.newDescription)
			if err != nil {
				t.Fatalf("Didn't expect an error updating task: %v", err)
			}

			if updatedTask.Title != tc.title.newTitle {
				t.Errorf("updated task title didn't match new title\nwant: %s, got: %s", tc.title.newTitle, updatedTask.Title)
			}
			
			if tc.description.newDescription != nil {
				if updatedTask.Description == nil {
					t.Errorf("Expected description to be updated, but got nil")
				} else if *updatedTask.Description != *tc.description.newDescription {
					t.Errorf("updated task description didn't match new description\nwant: %s, got: %s", *tc.description.newDescription, *updatedTask.Description)
				}
			}
		})
	}
}

func TestTaskClient_GetTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	title := "Test Task"
	description := "Test Description"

	createdTask, err := client.CreateTask(ctx, title, &description)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	retrievedTask, err := client.GetTask(ctx, createdTask.Id)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if retrievedTask.Id != createdTask.Id {
		t.Errorf("Expected ID %s, got %s", createdTask.Id, retrievedTask.Id)
	}
	if retrievedTask.Title != title {
		t.Errorf("Expected title %s, got %s", title, retrievedTask.Title)
	}
	if retrievedTask.Description == nil || *retrievedTask.Description != description {
		t.Errorf("Expected description %s, got %v", description, retrievedTask.Description)
	}
}

func TestTaskClient_GetTask_NonExistent(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	
	retrievedTask, err := client.GetTask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}

	if retrievedTask != nil {
		t.Errorf("Expected nil task for non-existent ID, got %v", retrievedTask)
	}
}

func TestTaskClient_DeleteTask(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	title := "Test Task"
	description := "Test Description"

	createdTask, err := client.CreateTask(ctx, title, &description)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	if len(mockClient.tasks) != 1 {
		t.Fatalf("Expected 1 task in mock, got %d", len(mockClient.tasks))
	}

	success, err := client.DeleteTask(ctx, createdTask.Id)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	if !success {
		t.Error("Expected DeleteTask to return true")
	}

	if len(mockClient.tasks) != 0 {
		t.Errorf("Expected 0 tasks in mock after deletion, got %d", len(mockClient.tasks))
	}

	retrievedTask, err := client.GetTask(ctx, createdTask.Id)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if retrievedTask != nil {
		t.Error("Expected task to be deleted, but it still exists")
	}
}

func TestTaskClient_DeleteTask_NonExistent(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks: make(map[string]*pb.Task),
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()

	success, err := client.DeleteTask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	if !success {
		t.Error("Expected DeleteTask to return true even for non-existent task")
	}
}

func TestTaskClient_CreateTask_Error(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks:     make(map[string]*pb.Task),
		shouldErr: true,
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()
	title := "Test Task"

	_, err := client.CreateTask(ctx, title, nil)
	if err == nil {
		t.Fatal("Expected error from CreateTask, got nil")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "failed to create task") {
		t.Errorf("Expected error to contain 'failed to create task', got: %v", err)
	}
}

func TestTaskClient_GetTask_Error(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks:     make(map[string]*pb.Task),
		shouldErr: true,
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()

	_, err := client.GetTask(ctx, "test-id")
	if err == nil {
		t.Fatal("Expected error from GetTask, got nil")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "failed to get task") {
		t.Errorf("Expected error to contain 'failed to get task', got: %v", err)
	}
}

func TestTaskClient_UpdateTask_Error(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks:     make(map[string]*pb.Task),
		shouldErr: true,
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()

	_, err := client.UpdateTask(ctx, "test-id", "new title", nil)
	if err == nil {
		t.Fatal("Expected error from UpdateTask, got nil")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "failed to update task") {
		t.Errorf("Expected error to contain 'failed to update task', got: %v", err)
	}
}

func TestTaskClient_DeleteTask_Error(t *testing.T) {
	mockClient := &mockTaskServiceClient{
		tasks:     make(map[string]*pb.Task),
		shouldErr: true,
	}

	client := &TaskClient{
		client: mockClient,
	}

	ctx := context.Background()

	_, err := client.DeleteTask(ctx, "test-id")
	if err == nil {
		t.Fatal("Expected error from DeleteTask, got nil")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "failed to delete task") {
		t.Errorf("Expected error to contain 'failed to delete task', got: %v", err)
	}
}
