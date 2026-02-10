package client

import (
	"context"
	"fmt"
	"net/http"

	taskv1 "github.com/jansdhillon/task-manager/proto/gen/task/v1"
	"github.com/jansdhillon/task-manager/proto/gen/task/v1/taskv1connect"
)

type TaskClient struct {
	client taskv1connect.TaskServiceClient
}

var NewTaskClient = newTaskClient

func newTaskClient(address string) (*TaskClient, error) {
	client := taskv1connect.NewTaskServiceClient(
		http.DefaultClient,
		address,
	)

	return &TaskClient{
		client: client,
	}, nil
}

func (c *TaskClient) CreateTask(ctx context.Context, title string, description *string) (*taskv1.Task, error) {
	req := &taskv1.CreateTaskRequest{
		Title: title,
	}
	if description != nil {
		req.Description = description
	}

	resp, err := c.client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return resp.Task, nil
}

func (c *TaskClient) GetTask(ctx context.Context, id string) (*taskv1.Task, error) {
	req := &taskv1.GetTaskRequest{
		Id: id,
	}

	resp, err := c.client.GetTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return resp.Task, nil
}

func (c *TaskClient) UpdateTask(ctx context.Context, id, title string, description *string, status *taskv1.Status) (*taskv1.Task, error) {
	req := &taskv1.UpdateTaskRequest{
		Id:    id,
		Title: title,
	}
	if description != nil {
		req.Description = description
	}
	if status != nil {
		req.Status = status
	}

	resp, err := c.client.UpdateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return resp.Task, nil
}

func (c *TaskClient) DeleteTask(ctx context.Context, id string) (bool, error) {
	req := &taskv1.DeleteTaskRequest{
		Id: id,
	}

	resp, err := c.client.DeleteTask(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to delete task: %w", err)
	}

	return resp.Success, nil
}
