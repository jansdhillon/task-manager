package main

import (
	"context"
	"fmt"
	"log"

	. "github.com/jansdhillon/task-client/internal"
	"github.com/urfave/cli/v3"
)

const (
	titleFlag       = "title"
	descriptionFlag = "description"
)

var createCmd = &cli.Command{
	Name:  "create",
	Usage: "Create a task with a title and, optionally, a description.",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() != 1 {
			return cli.Exit("Exactly one argument expected for address.", 1)
		}

		title := cmd.String(titleFlag)
		description := cmd.String(descriptionFlag)

		address := cmd.Args().Get(0)

		log.Printf("address: %s", address)
		log.Printf("title: %s", title)
		log.Printf("description: %T", description)

		result, err := ExecuteWithClient(address, func(c *TaskClient) (any, error) {
			var desc *string
			if description != "" {
				desc = &description
			}
			return c.CreateTask(ctx, title, desc)
		})

		if err != nil {
			errMsg := fmt.Sprintf("Error creating task: %v", err)
			return cli.Exit(errMsg, 1)
		}
		log.Printf("Result: %s", result)

		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     titleFlag,
			Usage:    "The title of the task to create.",
			Aliases:  []string{"t"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     descriptionFlag,
			Usage:    "(Optional) The description of the task to create.",
			Aliases:  []string{"d"},
			Required: false,
		},
	},
}
