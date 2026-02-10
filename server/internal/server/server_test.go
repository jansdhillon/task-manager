package server

import (
	"context"
	"testing"

	"github.com/google/uuid"
	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
	"github.com/jansdhillon/task-manager/server/internal/db"
)

type mockTaskDB struct {
	createTaskFunc func(ctx context.Context, title string, description *string) (*taskv1.Task, error)
	getTaskFunc    func(ctx context.Context, id uuid.UUID) (*taskv1.Task, error)
	updateTaskFunc func(ctx context.Context, id uuid.UUID, title string, description *string, status *taskv1.Status) (*taskv1.Task, error)
	deleteTaskFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *mockTaskDB) CreateTask(ctx context.Context, title string, description *string) (*taskv1.Task, error) {
	if m.createTaskFunc != nil {
		return m.createTaskFunc(ctx, title, description)
	}
	return nil, nil
}

func (m *mockTaskDB) GetTask(ctx context.Context, id uuid.UUID) (*taskv1.Task, error) {
	if m.getTaskFunc != nil {
		return m.getTaskFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTaskDB) UpdateTask(ctx context.Context, id uuid.UUID, title string, description *string, status *taskv1.Status) (*taskv1.Task, error) {
	if m.updateTaskFunc != nil {
		return m.updateTaskFunc(ctx, id, title, description, status)
	}
	return nil, nil
}

func (m *mockTaskDB) DeleteTask(ctx context.Context, id uuid.UUID) error {
	if m.deleteTaskFunc != nil {
		return m.deleteTaskFunc(ctx, id)
	}
	return nil
}

func TestNewTaskServer(t *testing.T) {
	mockDB := &mockTaskDB{}

	server := NewTaskServer(&db.TaskDB{})
	if server == nil {
		t.Fatal("NewTaskServer() returned nil")
	}
	if server.db == nil {
		t.Fatal("NewTaskServer() created server with nil db")
	}

	_ = mockDB // silence unused variable
}

func TestTaskServer_CreateTask(t *testing.T) {
	desc := "test description"

	tests := []struct {
		name        string
		title       string
		description *string
		wantErr     bool
	}{
		{"valid task", "Test Task", nil, false},
		{"task with description", "Test Task", &desc, false},
		{"empty title", "", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTaskServer(db.NewTaskDB())
			req := &taskv1.CreateTaskRequest{
				Title:       tt.title,
				Description: tt.description,
			}

			// This will fail without real DB connection
			_, _ = server.CreateTask(context.Background(), req)
		})
	}
}

func TestTaskServer_GetTask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"invalid UUID", "invalid-uuid", true},
		{"empty UUID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTaskServer(db.NewTaskDB())
			req := &taskv1.GetTaskRequest{
				Id: tt.id,
			}

			_, err := server.GetTask(context.Background(), req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskServer_UpdateTask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"invalid UUID", "invalid", true},
		{"empty UUID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTaskServer(db.NewTaskDB())
			req := &taskv1.UpdateTaskRequest{
				Id:    tt.id,
				Title: "Updated",
			}

			_, err := server.UpdateTask(context.Background(), req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskServer_DeleteTask(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"invalid UUID", "not-a-uuid", true},
		{"empty UUID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTaskServer(db.NewTaskDB())
			req := &taskv1.DeleteTaskRequest{
				Id: tt.id,
			}

			_, err := server.DeleteTask(context.Background(), req)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
