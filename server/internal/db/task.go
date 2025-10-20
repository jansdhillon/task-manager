package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jansdhillon/task-manager/server/.gen/postgres/public/model"
	. "github.com/jansdhillon/task-manager/server/.gen/postgres/public/table"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
	"github.com/jansdhillon/task-manager/server/internal/task"
)

type TaskDB struct {
	conn          *sql.DB
	secretsClient secrets.SecretsClient
}

func NewTaskDB(conn *sql.DB, sc secrets.SecretsClient) *TaskDB {
	return &TaskDB{
		conn:          conn,
		secretsClient: sc,
	}
}

func TaskFromDBTask(ctx context.Context, dbTask *model.Task) (*task.Task, error) {
	if dbTask == nil {
		return nil, errors.New("db task is nil")
	}

	status, err := statusFromModel(dbTask.Status)
	if err != nil {
		return nil, err
	}

	return &task.Task{
		ID:            dbTask.ID,
		Title:         dbTask.Title,
		Description:   dbTask.Description,
		CreatedAt:     dbTask.CreatedAt,
		LastUpdatedAt: dbTask.LastUpdatedAt,
		Deleted:       dbTask.Deleted,
		Status:        status,
	}, nil
}

func statusFromModel(status model.Status) (task.Status, error) {
	switch status {
	case model.Status_NotStarted:
		return task.NotStarted, nil
	case model.Status_InProgress:
		return task.InProgress, nil
	case model.Status_Completed:
		return task.Completed, nil
	default:
		return task.NotStarted, fmt.Errorf("unknown status %q", status)
	}
}

func (db *TaskDB) AddTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	if t == nil {
		return nil, errors.New("task is nil")
	}

	insertInput := model.Task{
		Title:       t.Title,
		Description: t.Description,
	}

	insertStmt := Task.
		INSERT(Task.Title, Task.Description).
		MODEL(insertInput).
		RETURNING(Task.AllColumns)

	var created model.Task

	if err := insertStmt.QueryContext(ctx, db.conn, &created); err != nil {
		return nil, err
	}

	return TaskFromDBTask(ctx, &created)
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
