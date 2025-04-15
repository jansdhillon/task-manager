package task

import (
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
	AddTask() (Task, error)
	DeleteTask() error
	UpdateTask() (Task, error)
	GetTask(id uuid.UUID) (Task, error)
}

func (t *Task) ToString() string {
	description := "(no description)"
	if t.Description != nil {
		description = *t.Description
	}
	return fmt.Sprintf("ID: %s, Title: %s, Description: %s, LastUpdatedAt: %s, CreatedAt: %s, Deleted: %t",
		t.ID, t.Title, description, t.LastUpdatedAt, t.CreatedAt, t.Deleted)
}

func New(title string, description string) *Task {
	return &Task{
		ID:            uuid.New(),
		Title:         title,
		Description:   &description,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Deleted:       false,
	}
}

type InMemoryTaskStore struct {
	Name  string
	Tasks []Task
}

func (s *InMemoryTaskStore) AddTask(t *Task) Task {
	s.Tasks = append(s.Tasks, *t)
	return *t
}

const MAX_TASKS = 10
