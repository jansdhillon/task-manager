package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
)

type mockDBConnector struct {
	connectFunc func(ctx context.Context) (*pgx.Conn, error)
}

func (m *mockDBConnector) Connect(ctx context.Context) (*pgx.Conn, error) {
	if m.connectFunc != nil {
		return m.connectFunc(ctx)
	}
	return nil, nil
}

type mockDBError struct {
	msg string
}

func (e *mockDBError) Error() string {
	return e.msg
}

func TestNewTaskDB(t *testing.T) {
	db := NewTaskDB()
	if db == nil {
		t.Fatal("NewTaskDB() returned nil")
	}
	if db.connector == nil {
		t.Fatal("NewTaskDB() created TaskDB with nil connector")
	}
}

func TestNewTaskDBWithConnector(t *testing.T) {
	mockConnector := &mockDBConnector{}
	db := NewTaskDBWithConnector(mockConnector)
	if db == nil {
		t.Fatal("NewTaskDBWithConnector() returned nil")
	}
	if db.connector != mockConnector {
		t.Error("NewTaskDBWithConnector() did not set the connector properly")
	}
}

func TestStatusConversions(t *testing.T) {
	tests := []struct {
		name     string
		proto    taskv1.Status
		expected taskv1.Status
	}{
		{"NOT_STARTED", taskv1.Status_STATUS_NOT_STARTED, taskv1.Status_STATUS_NOT_STARTED},
		{"IN_PROGRESS", taskv1.Status_STATUS_IN_PROGRESS, taskv1.Status_STATUS_IN_PROGRESS},
		{"COMPLETED", taskv1.Status_STATUS_COMPLETED, taskv1.Status_STATUS_COMPLETED},
		{"UNSPECIFIED", taskv1.Status_STATUS_UNSPECIFIED, taskv1.Status_STATUS_NOT_STARTED}, // UNSPECIFIED maps to NOT_STARTED
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := protoStatusToJetStatus(tt.proto)
			result := jetStatusToProtoStatus(model)

			if result != tt.expected {
				t.Errorf("statusProtoToModel->statusModelToProto = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTaskDB_WithConnectionError(t *testing.T) {
	mockConnector := &mockDBConnector{
		connectFunc: func(ctx context.Context) (*pgx.Conn, error) {
			return nil, &mockDBError{msg: "connection failed"}
		},
	}

	db := NewTaskDBWithConnector(mockConnector)
	ctx := context.Background()

	t.Run("CreateTask fails on connection error", func(t *testing.T) {
		_, err := db.CreateTask(ctx, "test", nil)
		if err == nil {
			t.Error("CreateTask() expected error but got none")
		}
	})

	t.Run("GetTask fails on connection error", func(t *testing.T) {
		_, err := db.GetTask(ctx, uuid.New())
		if err == nil {
			t.Error("GetTask() expected error but got none")
		}
	})

	t.Run("UpdateTask fails on connection error", func(t *testing.T) {
		_, err := db.UpdateTask(ctx, uuid.New(), "test", nil, nil)
		if err == nil {
			t.Error("UpdateTask() expected error but got none")
		}
	})

	t.Run("DeleteTask fails on connection error", func(t *testing.T) {
		err := db.DeleteTask(ctx, uuid.New())
		if err == nil {
			t.Error("DeleteTask() expected error but got none")
		}
	})
}

func TestTaskDB_WithCancelledContext(t *testing.T) {
	mockConnector := &mockDBConnector{
		connectFunc: func(ctx context.Context) (*pgx.Conn, error) {
			return nil, &mockDBError{msg: "context cancelled"}
		},
	}

	db := NewTaskDBWithConnector(mockConnector)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("CreateTask handles cancelled context", func(t *testing.T) {
		_, err := db.CreateTask(ctx, "test", nil)
		if err == nil {
			t.Error("CreateTask() expected error with cancelled context but got none")
		}
	})
}
