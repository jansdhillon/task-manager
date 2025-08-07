package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "task-client",
		Usage: "Connect to a Task Manager server to create and manage tasks.",
		Commands: []*cli.Command{
			connectCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
