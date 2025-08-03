package internal

import (
	"context"
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

	for i, tt := range tests {
		if tt.description != nil {
			t.Logf("\tTest: %d\tWhen using a title of %s and a description of %v\t", i, tt.title, *tt.description)
		} else {
			t.Logf("\tTest: %d\tWhen using a title of %s and a description of %v\t", i, tt.title, nil)
		}

		store := &InMemoryTaskStore{
			Name:  "Test Store",
			Tasks: make([]Task, 0),
		}

		server := NewTaskServer(store)

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
