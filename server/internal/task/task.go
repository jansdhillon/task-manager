package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	NotStarted Status = iota
	InProgress
	Completed
)

var statusName = map[Status]string{
	NotStarted: "not_started",
	InProgress: "in_progress",
	Completed:  "completed",
}

func (s Status) String() string {
	return statusName[s]
}

type Task struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Deleted       bool      `json:"deleted"`
	Status        Status    `json:"status"`
}

func (t *Task) String() string {
	description := "(no description)"
	if t.Description != nil {
		description = *t.Description
	}
	return fmt.Sprintf("ID: %s, Title: %s, Description: %s, LastUpdatedAt: %s, CreatedAt: %s, Deleted: %t",
		t.ID, t.Title, description, t.LastUpdatedAt, t.CreatedAt, t.Deleted)
}

func NewTask(title string, description *string) *Task {
	return &Task{
		ID:            uuid.New(),
		Title:         title,
		Description:   description,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Status:        NotStarted,
		Deleted:       false,
	}
}
