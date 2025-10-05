package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type InMemoryTaskStore struct {
	Name  string
	Tasks []*Task
}

func (s *InMemoryTaskStore) AddTask(ctx context.Context, t *Task) (*Task, error) {
	s.Tasks = append(s.Tasks, t)
	return t, nil
}

const MAX_TASKS = 10

// GetTask will search the store for a given
// Task by ID. If not found, an error will be
// returned.
func (s *InMemoryTaskStore) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	for i := range s.Tasks {
		if s.Tasks[i].ID == id {
			return s.Tasks[i], nil
		}
	}
	return nil, errors.New("task not found")
}

// DeleteTask takes the ID of a task and deletes it from
// the InMemoryTaskStore.
func (s *InMemoryTaskStore) DeleteTask(ctx context.Context, id uuid.UUID) error {
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
		return fmt.Errorf("task with ID %s not found", id)
	}

	s.Tasks = s.Tasks[:n]

	return nil
}

// UpdateTask updates the content of a task
func (s *InMemoryTaskStore) UpdateTask(ctx context.Context, id uuid.UUID, title string, description *string) (*Task, error) {
	task, err := s.GetTask(ctx, id)

	if err != nil {
		return nil, err
	}

	task.Title, task.Description = title, description

	return task, nil

}
