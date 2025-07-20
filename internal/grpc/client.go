package grpc

import (
	"context"
	"fmt"

	pb "github.com/jansdhillon/task-manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskClient struct {
	client pb.TaskServiceClient
	conn   *grpc.ClientConn
}

func NewTaskClient(address string) (*TaskClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := pb.NewTaskServiceClient(conn)
	return &TaskClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *TaskClient) Close() error {
	return c.conn.Close()
}

func (c *TaskClient) CreateTask(ctx context.Context, title string, description *string) (*pb.Task, error) {
	req := &pb.CreateTaskRequest{
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

func (c *TaskClient) GetTask(ctx context.Context, id string) (*pb.Task, error) {
	req := &pb.GetTaskRequest{
		Id: id,
	}

	resp, err := c.client.GetTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return resp.Task, nil
}

func (c *TaskClient) UpdateTask(ctx context.Context, id, title string, description *string) (*pb.Task, error) {
	req := &pb.UpdateTaskRequest{
		Id:    id,
		Title: title,
	}
	if description != nil {
		req.Description = description
	}

	resp, err := c.client.UpdateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return resp.Task, nil
}

func (c *TaskClient) DeleteTask(ctx context.Context, id string) (bool, error) {
	req := &pb.DeleteTaskRequest{
		Id: id,
	}

	resp, err := c.client.DeleteTask(ctx, req)
	if err != nil {
		return false, fmt.Errorf("failed to delete task: %w", err)
	}

	return resp.Success, nil
}