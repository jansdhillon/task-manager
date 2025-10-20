package task

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestAddTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}

	title := "world"
	var description *string

	addedTask, err := store.AddTask(t.Context(), title, description)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	if addedTask.Title != title {
		t.Errorf("got unexpected title: %s, want: %s", addedTask.Title, title)
	}

	if addedTask.Description != nil {
		t.Errorf("got unexpected title: %s, want: %v", *addedTask.Description, nil)
	}

}

// TestGetTask calls task.GetTask with an ID,
// checking for a valid return value.
func TestGetTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}
	desc := "world"
	title := "Hello World"

	task, err := store.AddTask(t.Context(), title, &desc)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	retrievedTask, err := store.GetTask(t.Context(), task.ID)

	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("Tasks didn't match! Original task: %v,\n retrieved task: %v", task, retrievedTask)
	}

}

// TestGetTaskBy still works if description
// is null.
func TestGetTaskNillable(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}
	ctx := context.Background()
	task, err := store.AddTask(ctx, "Hello World", nil)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	retrievedTask, err := store.GetTask(ctx, task.ID)

	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("Tasks didn't match! Original task: %v,\n retrieved task: %v", task, retrievedTask)
	}

}

func TestDeleteTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}

	var taskIDs []uuid.UUID
	ctx := context.Background()
	for i := range 10 {
		task, err := store.AddTask(ctx, fmt.Sprintf("task %d", i), nil)
		if err != nil {
			t.Error(err)
		}
		taskIDs = append(taskIDs, task.ID)

		returnedTask, err := store.GetTask(ctx, task.ID)
		if err != nil {
			t.Error(err)
		}

		if task.ID != returnedTask.ID {
			t.Error("ID mismatch")
		}
	}

	for _, id := range taskIDs {
		fmt.Printf("Should be deleting task with ID %s\n", id)
		err := store.DeleteTask(ctx, id)
		if err != nil {
			t.Error(err)
		}
	}

	if len(store.Tasks) != 0 {
		t.Errorf("expected all tasks to be deleted, but got %d left", len(store.Tasks))
	}
}

func TestDeleteTaskNotFoundRaises(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}

	for i := range 10 {
		task, err := store.AddTask(t.Context(), fmt.Sprintf("task %d", i), nil)
		if err != nil {
			t.Error(err)
		}

		_, err = store.GetTask(t.Context(), task.ID)
		if err != nil {
			t.Error(err)
		}
	}

	invalidID := uuid.New()
	fmt.Printf("Attempting to delete task with ID %s\n", invalidID)
	err := store.DeleteTask(t.Context(), invalidID)

	fmt.Printf("Received err: %v", err)

	if err == nil {
		t.Error("expected error")
	}

}

func TestUpdateTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}
	title := "task"
	var desc *string
	task, err := store.AddTask(t.Context(), title, desc)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	newDesc := "Hello"

	updatedTask, err := store.UpdateTask(t.Context(), task.ID, "New", &newDesc, InProgress)

	if err != nil {
		t.Error(err)
	}

	if updatedTask.Title != "New" && *updatedTask.Description != newDesc {
		t.Errorf("Title and desc didn't match! Updated title: %s, updated desc: %s", updatedTask.Title, *updatedTask.Description)
	}

	fmt.Printf("Updated task: %s, updated desc: %s\n", updatedTask.Title, *updatedTask.Description)

}

func TestUpdateTaskNotFoundRaises(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]*Task, 0)}
	ctx := context.Background()

	for i := range 10 {
		task, err := store.AddTask(t.Context(), fmt.Sprintf("task %d", i), nil)
		if err != nil {
			t.Error(err)
		}

		_, err = store.GetTask(ctx, task.ID)
		if err != nil {
			t.Error(err)
		}
	}

	invalidID := uuid.New()
	fmt.Printf("Attempting to update task with ID %s\n", invalidID)
	_, err := store.UpdateTask(ctx, invalidID, "h", nil, Completed)

	fmt.Printf("Received err: %v\n", err)

	if err == nil {
		t.Error("expected error")
	}

}
