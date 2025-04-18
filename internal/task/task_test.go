package task

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddTask(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}

	task := &Task{ID: uuid.New(), Title: "World", Description: nil, CreatedAt: time.Now(), LastUpdatedAt: time.Now(), Deleted: false}

	addedTask := store.AddTask(task)

	if addedTask != task {
		t.Errorf("Result was wrong, wanted: %s, got: %v.", task.String(), addedTask)
	}
}

// TestGetTaskById calls task.GetTaskById with an ID,
// checking for a valid return value.
func TestGetTaskById(t *testing.T) {
	store := &InMemoryTaskStore{Name: "Store", Tasks: make([]Task, 0)}
	task := New("Hello World", "world")
	fmt.Println(task.String())

	id := task.ID

	store.AddTask(task)

	retrievedTask, err := store.GetTaskById(id)

	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if !reflect.DeepEqual(retrievedTask, task) {
		t.Errorf("Tasks didn't match! Original task: %v,\n retrieved task: %v", task, retrievedTask)
	}

}
