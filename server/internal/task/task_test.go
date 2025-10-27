package task

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   string
	}{
		{
			name:   "not started status",
			status: NotStarted,
			want:   "not_started",
		},
		{
			name:   "in progress status",
			status: InProgress,
			want:   "in_progress",
		},
		{
			name:   "completed status",
			status: Completed,
			want:   "completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_String(t *testing.T) {
	id := uuid.New()
	createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	desc1 := "test description"
	tests := []struct {
		name string
		task Task
		want []string
	}{
		{
			name: "task with description",
			task: Task{
				ID:            id,
				Title:         "Test Task",
				Description:   &desc1,
				CreatedAt:     createdAt,
				LastUpdatedAt: updatedAt,
				Deleted:       false,
				Status:        NotStarted,
			},
			want: []string{id.String(), "Test Task", "test description", updatedAt.String(), createdAt.String(), "false"},
		},
		{
			name: "task without description",
			task: Task{
				ID:            id,
				Title:         "Test Task",
				Description:   nil,
				CreatedAt:     createdAt,
				LastUpdatedAt: updatedAt,
				Deleted:       true,
				Status:        InProgress,
			},
			want: []string{id.String(), "Test Task", "(no description)", updatedAt.String(), createdAt.String(), "true"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.task.String()
			for _, expectedPart := range tt.want {
				if !strings.Contains(got, expectedPart) {
					t.Errorf("Task.String() = %v, expected to contain %v", got, expectedPart)
				}
			}
		})
	}
}

func TestNewTask(t *testing.T) {
	desc1 := "test description"
	tests := []struct {
		name        string
		title       string
		description *string
		wantTitle   string
		wantDesc    *string
		wantDeleted bool
		wantStatus  Status
	}{
		{
			name:        "task with description",
			title:       "Test Task",
			description: &desc1,
			wantTitle:   "Test Task",
			wantDesc:    &desc1,
			wantDeleted: false,
			wantStatus:  NotStarted,
		},
		{
			name:        "task without description",
			title:       "Another Task",
			description: nil,
			wantTitle:   "Another Task",
			wantDesc:    nil,
			wantDeleted: false,
			wantStatus:  NotStarted,
		},
		{
			name:        "empty title task",
			title:       "",
			description: nil,
			wantTitle:   "",
			wantDesc:    nil,
			wantDeleted: false,
			wantStatus:  NotStarted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTask(tt.title, tt.description)
			
			if got == nil {
				t.Fatal("NewTask() returned nil")
			}
			
			if got.Title != tt.wantTitle {
				t.Errorf("NewTask() Title = %v, want %v", got.Title, tt.wantTitle)
			}
			
			if (got.Description == nil) != (tt.wantDesc == nil) {
				t.Errorf("NewTask() Description nil mismatch, got %v, want %v", got.Description == nil, tt.wantDesc == nil)
			}
			
			if got.Description != nil && tt.wantDesc != nil && *got.Description != *tt.wantDesc {
				t.Errorf("NewTask() Description = %v, want %v", *got.Description, *tt.wantDesc)
			}
			
			if got.Deleted != tt.wantDeleted {
				t.Errorf("NewTask() Deleted = %v, want %v", got.Deleted, tt.wantDeleted)
			}
			
			if got.Status != tt.wantStatus {
				t.Errorf("NewTask() Status = %v, want %v", got.Status, tt.wantStatus)
			}
			
			if got.ID == uuid.Nil {
				t.Error("NewTask() ID should not be nil UUID")
			}
			
			if got.CreatedAt.IsZero() {
				t.Error("NewTask() CreatedAt should not be zero")
			}
			
			if got.LastUpdatedAt.IsZero() {
				t.Error("NewTask() LastUpdatedAt should not be zero")
			}
			
			if got.CreatedAt.After(time.Now()) {
				t.Error("NewTask() CreatedAt should not be in the future")
			}
			
			if got.LastUpdatedAt.After(time.Now()) {
				t.Error("NewTask() LastUpdatedAt should not be in the future")
			}
		})
	}
}

func TestNewTask_UniqueIDs(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "generate 10 unique IDs",
			count: 10,
		},
		{
			name:  "generate 100 unique IDs",
			count: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := make(map[uuid.UUID]bool)
			
			for i := 0; i < tt.count; i++ {
				task := NewTask("test", nil)
				if ids[task.ID] {
					t.Errorf("NewTask() generated duplicate ID: %v", task.ID)
				}
				ids[task.ID] = true
			}
			
			if len(ids) != tt.count {
				t.Errorf("Expected %d unique IDs, got %d", tt.count, len(ids))
			}
		})
	}
}