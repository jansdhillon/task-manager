package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jansdhillon/task-manager/server/.gen/postgres/public/model"
	. "github.com/jansdhillon/task-manager/server/.gen/postgres/public/table"
	"github.com/jansdhillon/task-manager/server/internal/task"
)

type TaskDB struct {
	conn *sql.DB
}

func NewTaskDB(conn *sql.DB) *TaskDB {
	return &TaskDB{
		conn: conn,
	}
}

func (db *TaskDB) Close() error {
	db.conn.Close()

	return nil
}

func TaskFromDBTask(ctx context.Context, dbTask *model.Task) (*task.Task, error) {
	if dbTask == nil {
		return nil, errors.New("db task is nil")
	}

	status, err := taskStatusFromModelStatus(dbTask.Status)
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

func taskStatusFromModelStatus(status model.Status) (task.Status, error) {
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

func modelStatusFromTaskStatus(status task.Status) model.Status {
	switch status {
	case task.InProgress:
		return model.Status_InProgress
	case task.Completed:
		return model.Status_Completed
	default:
		return model.Status_NotStarted
	}
}

func (db *TaskDB) AddTask(ctx context.Context, title string, description *string) (*task.Task, error) {
	insertInput := model.Task{
		Title:       title,
		Description: description,
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

func (db *TaskDB) UpdateTask(ctx context.Context, id uuid.UUID, title string, description *string, status task.Status) (*task.Task, error) {
	original, err := db.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, fmt.Errorf("task not found")
	}

	in := model.Task{
		ID:          original.ID,
		Title:       title,
		Description: description,
		CreatedAt:   original.CreatedAt,
		Status:      modelStatusFromTaskStatus(status),
	}

	updateStmt := Task.UPDATE(Task.MutableColumns).MODEL(in).WHERE(Task.ID.EQ(postgres.UUID(id))).RETURNING(Task.AllColumns)

	var updatedTask model.Task

	if err := updateStmt.QueryContext(ctx, db.conn, &updatedTask); err != nil {
		return nil, err
	}

	return TaskFromDBTask(ctx, &updatedTask)
}

// Retrieve a task from the DB
func (db *TaskDB) GetTask(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	stmt := Task.SELECT(Task.AllColumns).WHERE(Task.ID.EQ(postgres.UUID(id)))

	var found model.Task
	if err := stmt.QueryContext(ctx, db.conn, &found); err != nil {
		return nil, err
	}

	return TaskFromDBTask(ctx, &found)
}

func (db *TaskDB) DeleteTask(ctx context.Context, id uuid.UUID) error {
	stmt := Task.DELETE().WHERE(Task.ID.EQ(postgres.UUID(id)))

	_, err := stmt.ExecContext(ctx, db.conn)
	if err != nil {
		return err
	}

	return nil
}

var _ task.TaskStore = &TaskDB{}
