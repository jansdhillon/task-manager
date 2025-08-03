package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}

	task := &Task{ID: uuid.New(), Title: "World", Description: nil, CreatedAt: time.Now(), LastUpdatedAt: time.Now(), Deleted: false}

	addedTask, err := store.AddTask(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	if addedTask.ID != task.ID {
		t.Errorf("Result was wrong, wanted: %s, got: %v.", task.String(), addedTask.String())
	}
}

// TestGetTask calls task.GetTask with an ID,
// checking for a valid return value.
func TestGetTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}
	desc := "world"
	task := NewTask("Hello World", &desc)
	fmt.Println(task.String())

	id := task.ID

	_, err := store.AddTask(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	retrievedTask, err := store.GetTask(id)

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
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}
	task := NewTask("Hello World", nil)
	fmt.Println(task.String())

	id := task.ID

	_, err := store.AddTask(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	retrievedTask, err := store.GetTask(id)

	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("Tasks didn't match! Original task: %v,\n retrieved task: %v", task, retrievedTask)
	}

}

func TestDeleteTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}

	var taskIDs []uuid.UUID

	for i := range 10 {
		task, err := store.AddTask(NewTask(fmt.Sprintf("task %d", i), nil))
		if err != nil {
			t.Error(err)
		}
		taskIDs = append(taskIDs, task.ID)

		returnedTask, err := store.GetTask(task.ID)
		if err != nil {
			t.Error(err)
		}

		if task.ID != returnedTask.ID {
			t.Error("ID mismatch")
		}
	}

	for _, id := range taskIDs {
		fmt.Printf("Should be deleting task with ID %s\n", id)
		err := store.DeleteTask(id)
		if err != nil {
			t.Error(err)
		}
	}

	if len(store.Tasks) != 0 {
		t.Errorf("expected all tasks to be deleted, but got %d left", len(store.Tasks))
	}
}

func TestDeleteTaskNotFoundRaises(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}

	var taskIDs []uuid.UUID

	for i := range 10 {
		task, err := store.AddTask(NewTask(fmt.Sprintf("task %d", i), nil))
		if err != nil {
			t.Error(err)
		}
		taskIDs = append(taskIDs, task.ID)

		_, err = store.GetTask(task.ID)
		if err != nil {
			t.Error(err)
		}
	}

	invalidID := uuid.New()
	fmt.Printf("Attempting to delete task with ID %s\n", invalidID)
	err := store.DeleteTask(invalidID)

	fmt.Printf("Received err: %v", err)

	if err == nil {
		t.Error("expected error")
	}

}

func TestUpdateTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}
	task := NewTask("task", nil)
	_, err := store.AddTask(task)
	if err != nil {
		t.Errorf("Failed to add task: %v", err)
	}

	newDesc := "Hello"

	updatedTask, err := store.UpdateTask(task.ID, "New", &newDesc)

	if err != nil {
		t.Error(err)
	}

	if updatedTask.Title != "New" && *updatedTask.Description != newDesc {
		t.Errorf("Title and desc didn't match! Updated title: %s, updated desc: %s", updatedTask.Title, *updatedTask.Description)
	}

	fmt.Printf("Updated task: %s, updated desc: %s\n", updatedTask.Title, *updatedTask.Description)

}

func TestUpdateTaskNotFoundRaises(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}

	var taskIDs []uuid.UUID

	for i := range 10 {
		task, err := store.AddTask(NewTask(fmt.Sprintf("task %d", i), nil))
		if err != nil {
			t.Error(err)
		}
		taskIDs = append(taskIDs, task.ID)

		_, err = store.GetTask(task.ID)
		if err != nil {
			t.Error(err)
		}
	}

	invalidID := uuid.New()
	fmt.Printf("Attempting to update task with ID %s\n", invalidID)
	_, err := store.UpdateTask(invalidID, "h", nil)

	fmt.Printf("Received err: %v\n", err)

	if err == nil {
		t.Error("expected error")
	}

}
