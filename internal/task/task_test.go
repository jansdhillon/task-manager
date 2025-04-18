package task

import (
	"fmt"
	"testing"
)

// TestGetTaskById calls task.GetTaskById with an ID,
// checking for a valid return value.
func TestGetTaskById(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}
	task := NewTask("Hello World", "world")
	fmt.Println(task.String())

	id := task.ID

	store.AddTask(task)

	task, err := store.GetTaskById(id)

	if err != nil {
		t.Error("Task not found!")
	}

}
