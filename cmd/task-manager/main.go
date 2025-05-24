// The main package for the Task Manager executable

package main

import (
	"context"
	"fmt"
	"os"

	task "github.com/jansdhillon/task-manager/internal/task"
	"github.com/urfave/cli/v3"
)

func main() {
	tasks := make([]task.Task, 0, task.MAX_TASKS)

	store := &task.InMemoryTaskStore{
		Name:  "Huel",
		Tasks: tasks,
	}

	for i := range task.MAX_TASKS {
		description := "world"
		task := task.New(fmt.Sprintf("Task #%d", i), &description)
		addedTask := store.AddTask(task)
		fmt.Printf("Added task: %s\n", addedTask.String())
	}

	(&cli.Command{}).Run(context.Background(), os.Args)

}
