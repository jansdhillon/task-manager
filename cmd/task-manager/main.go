// The main package for the Task Manager executable

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	task "github.com/jansdhillon/task-manager/internal/task"
	cli "github.com/urfave/cli/v3"
)

const (
	listFlag = "list"
)

func newApp() *cli.Command {
	tasks := make([]task.Task, 0, task.MAX_TASKS)

	store := &task.InMemoryTaskStore{
		Name:  "Huel",
		Tasks: tasks,
	}

	for i := range task.MAX_TASKS {
		description := "world"
		task := task.New(fmt.Sprintf("Task #%d", i), &description)
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

func actionSetup(c *cli.Command) (err error) {
	return nil
}

func main() {
	app := newApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
