// The main package for the Task Manager executable

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	. "github.com/jansdhillon/task-manager/internal"
	cli "github.com/urfave/cli/v3"
)

const (
	listFlag = "list"
)

func newApp() *cli.Command {
	tasks := make([]Task, 0, MAX_TASKS)

	store := &InMemoryTaskStore{
		Name:  "Huel",
		Tasks: tasks,
	}

	for i := range MAX_TASKS {
		description := "world"
		task := NewTask(fmt.Sprintf("Task #%d", i), &description)
		store.AddTask(task)
	}

	listCmd := &cli.Command{
		Name:  "tasks",
		Usage: "Show tasks in the store.",
		Action: func(context.Context, *cli.Command) error {
			for _, task := range store.Tasks {
				fmt.Printf("Task: %s\n", task.String())
			}
			return nil
		},
	}
	return &cli.Command{
		Usage: "Manage tasks.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    listFlag,
				Aliases: []string{"l"},
				Usage:   "List the tasks.",
			},
		},
		Commands: []*cli.Command{
			listCmd,
		},
	}
}

func main() {
	app := newApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
