package db

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
	"github.com/jansdhillon/task-manager/server/.gen/postgres/public/model"
	. "github.com/jansdhillon/task-manager/server/.gen/postgres/public/table"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DBConnector interface {
	Connect(ctx context.Context) (*pgx.Conn, error)
}

type ProductionDBConnector struct{}

func (p *ProductionDBConnector) Connect(ctx context.Context) (*pgx.Conn, error) {
	return Connect(ctx)
}

type TaskDB struct {
	connector DBConnector
}

func NewTaskDB() *TaskDB {
	return &TaskDB{
		connector: &ProductionDBConnector{},
	}
}

func NewTaskDBWithConnector(connector DBConnector) *TaskDB {
	return &TaskDB{
		connector: connector,
	}
}

func jetToProto(m *model.Task) *taskv1.Task {
	proto := &taskv1.Task{
		Id:            m.ID.String(),
		Title:         m.Title,
		CreatedAt:     timestamppb.New(m.CreatedAt),
		LastUpdatedAt: timestamppb.New(m.LastUpdatedAt),
		Deleted:       m.Deleted,
		Status:        jetStatusToProtoStatus(m.Status),
	}
	if m.Description != nil {
		proto.Description = m.Description
	}
	return proto
}

func protoStatusToJetStatus(s taskv1.Status) model.Status {
	switch s {
	case taskv1.Status_STATUS_NOT_STARTED:
		return model.Status_NotStarted
	case taskv1.Status_STATUS_IN_PROGRESS:
		return model.Status_InProgress
	case taskv1.Status_STATUS_COMPLETED:
		return model.Status_Completed
	case taskv1.Status_STATUS_UNSPECIFIED:
		return model.Status_NotStarted
	default:
		return model.Status_NotStarted
	}
}

func jetStatusToProtoStatus(s model.Status) taskv1.Status {
	switch s {
	case model.Status_NotStarted:
		return taskv1.Status_STATUS_NOT_STARTED
	case model.Status_InProgress:
		return taskv1.Status_STATUS_IN_PROGRESS
	case model.Status_Completed:
		return taskv1.Status_STATUS_COMPLETED
	default:
		return taskv1.Status_STATUS_UNSPECIFIED
	}
}

func (t *TaskDB) CreateTask(ctx context.Context, title string, description *string) (*taskv1.Task, error) {
	conn, err := t.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	newTask := model.Task{
		ID:     uuid.New(),
		Title:  title,
		Status: model.Status_NotStarted,
	}
	if description != nil {
		newTask.Description = description
	}

	stmt := Task.INSERT(Task.ID, Task.Title, Task.Description, Task.Status).
		MODEL(newTask).
		RETURNING(Task.AllColumns)

	var inserted model.Task
	err = stmt.QueryContext(ctx, db, &inserted)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	return jetToProto(&inserted), nil
}

func (t *TaskDB) GetTask(ctx context.Context, id uuid.UUID) (*taskv1.Task, error) {
	stmt := postgres.SELECT(Task.AllColumns).
		FROM(Task).
		WHERE(Task.ID.EQ(postgres.UUID(id)))

	conn, err := t.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	var dbTask model.Task
	err = stmt.QueryContext(ctx, db, &dbTask)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return jetToProto(&dbTask), nil
}

func (t *TaskDB) UpdateTask(ctx context.Context, id uuid.UUID, title string, description *string, status *taskv1.Status) (*taskv1.Task, error) {
	conn, err := t.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	updateStmt := Task.UPDATE(Task.Title, Task.Description)
	if status != nil {
		updateStmt = Task.UPDATE(Task.Title, Task.Description, Task.Status)
	}

	values := map[postgres.Column]any{
		Task.Title:       title,
		Task.Description: description,
	}
	if status != nil {
		values[Task.Status] = protoStatusToJetStatus(*status)
	}

	stmt := updateStmt.
		SET(values).
		WHERE(Task.ID.EQ(postgres.UUID(id))).
		RETURNING(Task.AllColumns)

	var updated model.Task
	err = stmt.QueryContext(ctx, db, &updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return jetToProto(&updated), nil
}

func (t *TaskDB) DeleteTask(ctx context.Context, id uuid.UUID) error {
	conn, err := t.connector.Connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	db := stdlib.OpenDB(*conn.Config())
	defer db.Close()

	stmt := Task.UPDATE(Task.Deleted).
		SET(true).
		WHERE(Task.ID.EQ(postgres.UUID(id)))

	_, err = stmt.ExecContext(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
