package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
	"github.com/jansdhillon/task-manager/server/internal/task"
)

type TaskDB struct {
	store         task.TaskStore
	secretsClient secrets.SecretsClient
}

func (db *TaskDB) AddTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	return &task.Task{}, nil
}

func (db *TaskDB) UpdateTask(ctx context.Context, id uuid.UUID, title string, description *string) (*task.Task, error) {
	return &task.Task{}, nil
}

func (db *TaskDB) GetTask(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	return &task.Task{}, nil
}

func (db *TaskDB) DeleteTask(ctx context.Context, uuid uuid.UUID) error {
	return nil
}

var _ task.TaskStore = &TaskDB{}
