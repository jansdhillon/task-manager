package task

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Deleted       bool      `json:"deleted"`
}

type TaskStore interface {
	AddTask(t *Task) (Task, error)
	DeleteTask(id uuid.UUID) error
	UpdateTask(id uuid.UUID) (Task, error)
	GetTask(id uuid.UUID) (Task, error)
}

func (t *Task) String() string {
	description := "(no description)"
	if t.Description != nil {
		description = *t.Description
	}
	return fmt.Sprintf("ID: %s, Title: %s, Description: %s, LastUpdatedAt: %s, CreatedAt: %s, Deleted: %t",
		t.ID, t.Title, description, t.LastUpdatedAt, t.CreatedAt, t.Deleted)
}

func New(title string, description *string) *Task {
	return &Task{
		ID:            uuid.New(),
		Title:         title,
		Description:   description,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Deleted:       false,
	}
}

type InMemoryTaskStore struct {
	Name  string
	Tasks []Task
}

func (s *InMemoryTaskStore) AddTask(t *Task) *Task {
	s.Tasks = append(s.Tasks, *t)
	return t
}

const MAX_TASKS = 10

// GetTask will search the store for a given
// Task by ID. If not found, an error will be
// returned.
func (s *InMemoryTaskStore) GetTask(id uuid.UUID) (*Task, error) {
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return &s.Tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

// DeleteTask takes the ID of a task and deletes it from
// the InMemoryTaskStore.
func (s *InMemoryTaskStore) DeleteTask(id uuid.UUID) error {
	n := 0
	found := false
	for _, task := range s.Tasks {
		if task.ID != id {
			s.Tasks[n] = task
			n++
		} else {
			found = true
			fmt.Printf("Deleting task with ID %s\n", task.ID)
		}
	}

	if !found {
		return errors.New(fmt.Sprintf("Task with ID %s not found!\n", id))
	}

	s.Tasks = s.Tasks[:n]

	return nil
}

// UpdateTask updates the content of a task
func (s *InMemoryTaskStore) UpdateTask(id uuid.UUID, title string, description *string) (*Task, error) {
	task, err := s.GetTask(id)

	if err != nil {
		return nil, err
	}

	task.Title, task.Description = title, description

	return task, nil

}
