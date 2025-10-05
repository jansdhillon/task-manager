package server

import (
	"context"
	"net"
	"strings"
	"testing"

	"github.com/google/uuid"
	pb "github.com/jansdhillon/task-manager/proto"
	"github.com/jansdhillon/task-manager/server/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func setupTestServer(t *testing.T) (*TaskServer, pb.TaskServiceClient, func()) {
	t.Helper()

	store := &task.InMemoryTaskStore{
		Name:  "Test Store",
		Tasks: make([]*task.Task, 0),
	}

	taskServer := NewTaskServer(store)

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, taskServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewTaskServiceClient(conn)

	cleanup := func() {
		conn.Close()
		s.Stop()
	}

	return taskServer, client, cleanup
}

func stringPtr(s string) *string {
	return &s
}

func TestTaskServer_CreateTask(t *testing.T) {
	tests := []struct {
		name        string
		request     *pb.CreateTaskRequest
		wantErr     bool
		wantErrCode codes.Code
		validate    func(t *testing.T, resp *pb.CreateTaskResponse)
	}{
		{
			name: "create task with title only",
			request: &pb.CreateTaskRequest{
				Title: "Test Task",
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.CreateTaskResponse) {
				if resp.Task.Title != "Test Task" {
					t.Errorf("Expected title 'Test Task', got %q", resp.Task.Title)
				}
				if resp.Task.Description != nil {
					t.Errorf("Expected nil description, got %v", resp.Task.Description)
				}
				if resp.Task.Id == "" {
					t.Error("Expected non-empty task ID")
				}
				if resp.Task.CreatedAt == nil {
					t.Error("Expected non-nil CreatedAt")
				}
				if resp.Task.LastUpdatedAt == nil {
					t.Error("Expected non-nil LastUpdatedAt")
				}
				if resp.Task.Deleted {
					t.Error("Expected task to not be deleted")
				}
			},
		},
		{
			name: "create task with title and description",
			request: &pb.CreateTaskRequest{
				Title:       "Task with Description",
				Description: stringPtr("This is a test description"),
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.CreateTaskResponse) {
				if resp.Task.Title != "Task with Description" {
					t.Errorf("Expected title 'Task with Description', got %q", resp.Task.Title)
				}
				if resp.Task.Description == nil {
					t.Error("Expected non-nil description")
				} else if *resp.Task.Description != "This is a test description" {
					t.Errorf("Expected description 'This is a test description', got %q", *resp.Task.Description)
				}
			},
		},
		{
			name: "create task with empty description",
			request: &pb.CreateTaskRequest{
				Title:       "Task with Empty Description",
				Description: stringPtr(""),
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.CreateTaskResponse) {
				if resp.Task.Description == nil {
					t.Error("Expected non-nil description")
				} else if *resp.Task.Description != "" {
					t.Errorf("Expected empty description, got %q", *resp.Task.Description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, client, cleanup := setupTestServer(t)
			defer cleanup()

			resp, err := client.CreateTask(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.wantErrCode != codes.OK {
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %T", err)
						return
					}
					if st.Code() != tt.wantErrCode {
						t.Errorf("Expected error code %v, got %v", tt.wantErrCode, st.Code())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, resp)
			}
		})
	}
}

func TestTaskServer_UpdateTask(t *testing.T) {
	tests := []struct {
		name        string
		setupTask   func(t *testing.T, server *TaskServer) string
		request     *pb.UpdateTaskRequest
		wantErr     bool
		wantErrCode codes.Code
		validate    func(t *testing.T, resp *pb.UpdateTaskResponse, originalTaskID string)
	}{
		{
			name: "update existing task title and description",
			setupTask: func(t *testing.T, server *TaskServer) string {
				resp, err := server.CreateTask(context.Background(), &pb.CreateTaskRequest{
					Title:       "Original Title",
					Description: stringPtr("Original Description"),
				})
				if err != nil {
					t.Fatalf("Failed to create test task: %v", err)
				}
				return resp.Task.Id
			},
			request: &pb.UpdateTaskRequest{
				Title:       "Updated Title",
				Description: stringPtr("Updated Description"),
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.UpdateTaskResponse, originalTaskID string) {
				if resp.Task.Id != originalTaskID {
					t.Errorf("Expected task ID %q, got %q", originalTaskID, resp.Task.Id)
				}
				if resp.Task.Title != "Updated Title" {
					t.Errorf("Expected title 'Updated Title', got %q", resp.Task.Title)
				}
				if resp.Task.Description == nil {
					t.Error("Expected non-nil description")
				} else if *resp.Task.Description != "Updated Description" {
					t.Errorf("Expected description 'Updated Description', got %q", *resp.Task.Description)
				}
			},
		},
		{
			name: "update task with nil description",
			setupTask: func(t *testing.T, server *TaskServer) string {
				resp, err := server.CreateTask(context.Background(), &pb.CreateTaskRequest{
					Title:       "Original Title",
					Description: stringPtr("Original Description"),
				})
				if err != nil {
					t.Fatalf("Failed to create test task: %v", err)
				}
				return resp.Task.Id
			},
			request: &pb.UpdateTaskRequest{
				Title:       "Updated Title",
				Description: nil,
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.UpdateTaskResponse, originalTaskID string) {
				if resp.Task.Description != nil {
					t.Errorf("Expected nil description, got %q", *resp.Task.Description)
				}
			},
		},
		{
			name: "update non-existent task",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return uuid.New().String()
			},
			request: &pb.UpdateTaskRequest{
				Title: "Updated Title",
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
		{
			name: "update task with invalid UUID",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return "invalid-uuid"
			},
			request: &pb.UpdateTaskRequest{
				Title: "Updated Title",
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client, cleanup := setupTestServer(t)
			defer cleanup()

			taskID := tt.setupTask(t, server)
			tt.request.Id = taskID

			resp, err := client.UpdateTask(context.Background(), tt.request)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.wantErrCode != codes.OK {
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %T", err)
						return
					}
					if st.Code() != tt.wantErrCode {
						t.Errorf("Expected error code %v, got %v", tt.wantErrCode, st.Code())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, resp, taskID)
			}
		})
	}
}

func TestTaskServer_GetTask(t *testing.T) {
	tests := []struct {
		name        string
		setupTask   func(t *testing.T, server *TaskServer) string
		wantErr     bool
		wantErrCode codes.Code
		validate    func(t *testing.T, resp *pb.GetTaskResponse, expectedTaskID string)
	}{
		{
			name: "get existing task",
			setupTask: func(t *testing.T, server *TaskServer) string {
				resp, err := server.CreateTask(context.Background(), &pb.CreateTaskRequest{
					Title:       "Test Task",
					Description: stringPtr("Test Description"),
				})
				if err != nil {
					t.Fatalf("Failed to create test task: %v", err)
				}
				return resp.Task.Id
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.GetTaskResponse, expectedTaskID string) {
				if resp.Task.Id != expectedTaskID {
					t.Errorf("Expected task ID %q, got %q", expectedTaskID, resp.Task.Id)
				}
				if resp.Task.Title != "Test Task" {
					t.Errorf("Expected title 'Test Task', got %q", resp.Task.Title)
				}
				if resp.Task.Description == nil {
					t.Error("Expected non-nil description")
				} else if *resp.Task.Description != "Test Description" {
					t.Errorf("Expected description 'Test Description', got %q", *resp.Task.Description)
				}
				if resp.Task.CreatedAt == nil {
					t.Error("Expected non-nil CreatedAt")
				}
				if resp.Task.LastUpdatedAt == nil {
					t.Error("Expected non-nil LastUpdatedAt")
				}
				if resp.Task.Deleted {
					t.Error("Expected task to not be deleted")
				}
			},
		},
		{
			name: "get non-existent task",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return uuid.New().String()
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
		{
			name: "get task with invalid UUID",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return "invalid-uuid"
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client, cleanup := setupTestServer(t)
			defer cleanup()

			taskID := tt.setupTask(t, server)

			resp, err := client.GetTask(context.Background(), &pb.GetTaskRequest{
				Id: taskID,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.wantErrCode != codes.OK {
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %T", err)
						return
					}
					if st.Code() != tt.wantErrCode {
						t.Errorf("Expected error code %v, got %v", tt.wantErrCode, st.Code())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, resp, taskID)
			}
		})
	}
}

func TestTaskServer_DeleteTask(t *testing.T) {
	tests := []struct {
		name        string
		setupTask   func(t *testing.T, server *TaskServer) string
		wantErr     bool
		wantErrCode codes.Code
		validate    func(t *testing.T, resp *pb.DeleteTaskResponse, taskID string, server *TaskServer)
	}{
		{
			name: "delete existing task",
			setupTask: func(t *testing.T, server *TaskServer) string {
				resp, err := server.CreateTask(context.Background(), &pb.CreateTaskRequest{
					Title:       "Task to Delete",
					Description: stringPtr("This task will be deleted"),
				})
				if err != nil {
					t.Fatalf("Failed to create test task: %v", err)
				}
				return resp.Task.Id
			},
			wantErr: false,
			validate: func(t *testing.T, resp *pb.DeleteTaskResponse, taskID string, server *TaskServer) {
				if !resp.Success {
					t.Error("Expected successful deletion")
				}
				_, err := server.GetTask(context.Background(), &pb.GetTaskRequest{Id: taskID})
				if err == nil {
					t.Error("Expected task to be deleted and not retrievable")
				}
				if !strings.Contains(err.Error(), "task not found") {
					t.Errorf("Expected 'task not found' error, got: %v", err)
				}
			},
		},
		{
			name: "delete non-existent task",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return uuid.New().String()
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
		{
			name: "delete task with invalid UUID",
			setupTask: func(t *testing.T, server *TaskServer) string {
				return "invalid-uuid"
			},
			wantErr:     true,
			wantErrCode: codes.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client, cleanup := setupTestServer(t)
			defer cleanup()

			taskID := tt.setupTask(t, server)

			resp, err := client.DeleteTask(context.Background(), &pb.DeleteTaskRequest{
				Id: taskID,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
					return
				}
				if tt.wantErrCode != codes.OK {
					st, ok := status.FromError(err)
					if !ok {
						t.Errorf("Expected gRPC status error, got %T", err)
						return
					}
					if st.Code() != tt.wantErrCode {
						t.Errorf("Expected error code %v, got %v", tt.wantErrCode, st.Code())
					}
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, resp, taskID, server)
			}
		})
	}
}
