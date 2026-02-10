package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jansdhillon/task-manager/client/internal/client"
	"github.com/urfave/cli/v3"
)

var createCmd = &cli.Command{
	Name:  "create",
	Usage: "Create a task with a title and, optionally, a description.",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() == 0 || cmd.NArg() > 2 {
			return cli.Exit("Expected one or two arguments.", 1)
		}

		address := cmd.String(serverAddressFlag)
		title := cmd.Args().Get(0)
		description := cmd.Args().Get(1)
		var desc *string
		if description != "" {
			desc = &description
		}

		c, err := client.NewTaskClient(address)
		if err != nil {
			errMsg := fmt.Sprintf("Error creating client: %v", err)
			return cli.Exit(errMsg, 1)
		}

		task, err := c.CreateTask(ctx, title, desc)
		if err != nil {
			errMsg := fmt.Sprintf("Error creating task: %v", err)
			return cli.Exit(errMsg, 1)
		}
		log.Printf("Created task: %s", task)

		return nil
	},
}
