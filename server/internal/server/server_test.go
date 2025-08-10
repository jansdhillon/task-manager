package server

import (
	"context"
	"testing"

	pb "github.com/jansdhillon/task-manager/proto"
	"github.com/jansdhillon/task-manager/server/internal/task"
)

// TestCreateTaskServer tests the behavior of TaskServer.CreateTask.
// It verifies that tasks can be created with both nil and non-nil descriptions,
// and that all required fields are properly set on the created task.
func TestCreateTaskServer(t *testing.T) {
	notNullDesc := "world"
	tests := []struct {
		title       string
		description *string
	}{
		{"hello", &notNullDesc},
		{"hello", nil},
	}
	t.Log("Given the need to test creating a task from a TaskRequest")

	store := &task.InMemoryTaskStore{
		Name:  "Test Store",
		Tasks: make([]task.Task, 0),
	}

	server := NewTaskServer(store)

	for i, tt := range tests {
		if tt.description != nil {
			t.Logf("\tTest: %d\tWhen using a title of %s and a description of %v\t", i, tt.title, *tt.description)
		} else {
			t.Logf("\tTest: %d\tWhen using a title of %s and a description of %v\t", i, tt.title, nil)
		}

		req := &pb.CreateTaskRequest{
			Title:       tt.title,
			Description: tt.description,
		}

		resp, err := server.CreateTask(context.Background(), req)
		if err != nil {
			t.Fatalf("\t\tShould be able to create task but got error: %v", err)
		}

		if resp.Task.Title != tt.title {
			t.Errorf("\t\tExpected title %q but got %q", tt.title, resp.Task.Title)
		}

		if tt.description != nil {
			if resp.Task.Description == nil {
				t.Errorf("\t\tExpected description %q but got nil", *tt.description)
			} else if *resp.Task.Description != *tt.description {
				t.Errorf("\t\tExpected description %q but got %q", *tt.description, *resp.Task.Description)
			}
		} else {
			if resp.Task.Description != nil {
				t.Errorf("\t\tExpected nil description but got %q", *resp.Task.Description)
			}
		}

		if resp.Task.Id == "" {
			t.Error("\t\tExpected task to have an ID")
		}

		if resp.Task.CreatedAt == nil {
			t.Error("\t\tExpected task to have a creation timestamp")
		}

		if resp.Task.LastUpdatedAt == nil {
			t.Error("\t\tExpected task to have a last updated timestamp")
		}

		if resp.Task.Deleted {
			t.Error("\t\tExpected task to not be deleted")
		}

		t.Logf("\t\tShould create task successfully")
	}
}

// TestUpdateTaskServer tests the behavior of TaskServer.UpdateTask.
// It verifies that existing tasks can be updated with different description scenarios
// (with description, empty description, nil description) and that timestamps are updated.
func TestUpdateTaskServer(t *testing.T) {
	newTitle := "updated title"
	newDescription := "updated description"
	emptyDescription := ""

	tests := []struct {
		name        string
		title       string
		description *string
	}{
		{"update with description", newTitle, &newDescription},
		{"update with empty description", newTitle, &emptyDescription},
		{"update with nil description", newTitle, nil},
	}

	t.Log("Given the need to test updating a task")

	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen updating task with %s", i, tt.name)

		store := &task.InMemoryTaskStore{
			Name:  "Test Store",
			Tasks: make([]task.Task, 0),
		}

		server := NewTaskServer(store)

		originalDesc := "original description"
		originalTask := task.NewTask("Test Task", &originalDesc)
		createdTask, _ := store.AddTask(originalTask)

		req := &pb.UpdateTaskRequest{
			Id:          createdTask.ID.String(),
			Title:       tt.title,
			Description: tt.description,
		}

		resp, err := server.UpdateTask(context.Background(), req)
		if err != nil {
			t.Fatalf("\t\tShould be able to update task but got error: %v", err)
		}

		if resp.Task.Title != tt.title {
			t.Errorf("\t\tExpected title %q but got %q", tt.title, resp.Task.Title)
		}

		if tt.description != nil {
			if resp.Task.Description == nil {
				t.Errorf("\t\tExpected description %q but got nil", *tt.description)
			} else if *resp.Task.Description != *tt.description {
				t.Errorf("\t\tExpected description %q but got %q", *tt.description, *resp.Task.Description)
			}
		} else {
			if resp.Task.Description != nil {
				t.Errorf("\t\tExpected nil description but got %q", *resp.Task.Description)
			}
		}

		if resp.Task.Id != createdTask.ID.String() {
			t.Errorf("\t\tExpected task ID to remain %q but got %q", createdTask.ID.String(), resp.Task.Id)
		}

		if resp.Task.LastUpdatedAt == nil {
			t.Error("\t\tExpected task to have updated timestamp")
		}

		t.Logf("\t\tShould update task successfully")
	}
}

// TestGetTaskServer tests the behavior of TaskServer.GetTask.
// It verifies that existing tasks can be retrieved successfully and that
// appropriate errors are returned for non-existent tasks.
func TestGetTaskServer(t *testing.T) {
	tests := []struct {
		name        string
		setupTask   bool
		expectError bool
	}{
		{"get existing task", true, false},
		{"get non-existent task", false, true},
	}

	t.Log("Given the need to test getting a task")

	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen %s", i, tt.name)

		store := &task.InMemoryTaskStore{
			Name:  "Test Store",
			Tasks: make([]task.Task, 0),
		}

		server := NewTaskServer(store)

		var taskID string
		if tt.setupTask {
			desc := "test description"
			task := task.NewTask("Test Task", &desc)
			createdTask, _ := store.AddTask(task)
			taskID = createdTask.ID.String()
		} else {
			taskID = "non-existent-id"
		}

		req := &pb.GetTaskRequest{
			Id: taskID,
		}

		resp, err := server.GetTask(context.Background(), req)

		if tt.expectError {
			if err == nil {
				t.Fatalf("\t\tShould receive error for non-existent task but got nil")
			}
			t.Logf("\t\tShould return error for non-existent task")
		} else {
			if err != nil {
				t.Fatalf("\t\tShould be able to get task but got error: %v", err)
			}

			if resp.Task.Id != taskID {
				t.Errorf("\t\tExpected task ID %q but got %q", taskID, resp.Task.Id)
			}

			if resp.Task.Title != "Test Task" {
				t.Errorf("\t\tExpected task title %q but got %q", "Test Task", resp.Task.Title)
			}

			if resp.Task.Description == nil || *resp.Task.Description != "test description" {
				t.Errorf("\t\tExpected task description %q but got %v", "test description", resp.Task.Description)
			}

			if resp.Task.CreatedAt == nil {
				t.Error("\t\tExpected task to have a creation timestamp")
			}

			if resp.Task.LastUpdatedAt == nil {
				t.Error("\t\tExpected task to have a last updated timestamp")
			}

			if resp.Task.Deleted {
				t.Error("\t\tExpected task to not be deleted")
			}

			t.Logf("\t\tShould get task successfully")
		}
	}
}

// TestDeleteTaskServer tests the behavior of TaskServer.DeleteTask.
// It verifies that existing tasks can be deleted successfully and that
// appropriate errors are returned for non-existent tasks.
func TestDeleteTaskServer(t *testing.T) {
	tests := []struct {
		name        string
		setupTask   bool
		expectError bool
	}{
		{"delete existing task", true, false},
		{"delete non-existent task", false, true},
	}

	t.Log("Given the need to test deleting a task")

	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen %s", i, tt.name)

		store := &task.InMemoryTaskStore{
			Name:  "Test Store",
			Tasks: make([]task.Task, 0),
		}

		server := NewTaskServer(store)

		var taskID string
		if tt.setupTask {
			desc := "test description"
			task := task.NewTask("Test Task", &desc)
			createdTask, _ := store.AddTask(task)
			taskID = createdTask.ID.String()
		} else {
			taskID = "non-existent-id"
		}

		req := &pb.DeleteTaskRequest{
			Id: taskID,
		}

		resp, err := server.DeleteTask(context.Background(), req)

		if tt.expectError {
			if err == nil {
				t.Fatalf("\t\tShould receive error for non-existent task but got nil")
			}
			t.Logf("\t\tShould return error for non-existent task")
		} else {
			if err != nil {
				t.Fatalf("\t\tShould be able to delete task but got error: %v", err)
			}

			if !resp.Success {
				t.Error("\t\tExpected successful deletion")
			}

			getReq := &pb.GetTaskRequest{Id: taskID}
			_, getErr := server.GetTask(context.Background(), getReq)
			if getErr == nil {
				t.Error("\t\tExpected task to be deleted and not retrievable")
			}

			t.Logf("\t\tShould delete task successfully")
		}
	}
}
