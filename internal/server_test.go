package internal

import (
	"context"
	"strconv"
	"testing"

	pb "github.com/jansdhillon/task-manager/proto"
)

func TestCreateTask(t *testing.T) {
	notNullDesc := "world"
	tests := []struct {
		title       string
		description *string
	}{
		{"hello", &notNullDesc},
		{"hello", nil},
	}
	t.Log("Given the need to test creating a task from a TaskRequest")

	store := &InMemoryTaskStore{
		Name:  "Test Store",
		Tasks: make([]Task, 0),
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

// TestUpdateTaskServer tests the behavior of
// TaskServer.UpdateTeask
func TestUpdateTaskServer(t *testing.T) {
	newTitle := "new title"
	newDescription := "new description"
	var nilDescription string

	tt := []struct {
		title       string
		description *string
	}{
		{newTitle, &newDescription},
		{newTitle, &nilDescription},
	}

	store := &InMemoryTaskStore{
		Name:  "Test Store",
		Tasks: make([]Task, 0),
	}

	server := NewTaskServer(store)

	for i, tc := range tt {
		desc := "original description"
		originalTask := NewTask(strconv.Itoa(i), &desc)
		createdTask, _ := store.AddTask(originalTask)

		req := &pb.UpdateTaskRequest{
			Id:          createdTask.ID.String(),
			Title:       tc.title,
			Description: tc.description,
		}

		res, err := server.UpdateTask(context.Background(), req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		updatedTask := res.Task

		if updatedTask.Title == originalTask.Title || (updatedTask.Description != nil && updatedTask.Description == originalTask.Description) {
			t.Errorf("Task was not updated!")
		}

	}
}
