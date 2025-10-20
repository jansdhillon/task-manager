package client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/jansdhillon/task-manager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskClient struct {
	client pb.TaskServiceClient
	conn   *grpc.ClientConn
}

var NewTaskClient = newTaskClient

func newTaskClient(address string, useInsecure bool) (*TaskClient, error) {
	var creds credentials.TransportCredentials
	if useInsecure {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(nil)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
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
	if c.conn == nil {
		return nil
	}
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

func ExecuteWithClient(address string, useInsecure bool, fn func(*TaskClient) (any, error)) (any, error) {
	c, err := NewTaskClient(address, useInsecure)
	if err != nil {
		log.Printf("error creating client: %v", err)
		return nil, err
	}

	defer func() {
		if cerr := c.Close(); cerr != nil {
			log.Printf("error closing client connection: %v", cerr)
		}
	}()

	res, err := fn(c)

	return res, err
}
