package db

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	. "github.com/jansdhillon/task-manager/server/.gen/postgres/public/table"

	"github.com/jansdhillon/task-manager/server/.gen/postgres/public/model"
)

type DBConnector interface {
	Connect() (*pgx.Conn, error)
}

type ProductionDBConnector struct{}

func (p *ProductionDBConnector) Connect() (*pgx.Conn, error) {
	return Connect()
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

func (t *TaskDB) GetTasks(ctx context.Context) ([]model.Task, error) {
	stmt := SELECT(
		Task.ID,
		Task.CreatedAt,
	).FROM(Task).ORDER_BY(Task.CreatedAt.ASC())

	conn, err := t.connector.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	db := stdlib.OpenDB(*conn.Config())

	var tasks []model.Task
	err = stmt.QueryContext(ctx, db, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
