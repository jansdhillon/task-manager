package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

const (
	serverAddressFlag = "address"
	insecureFlag      = "insecure"
)

func main() {
	cmd := &cli.Command{
		Name:  "task-client",
		Usage: "Connect to a Task Manager server to create and manage tasks.",
		Commands: []*cli.Command{
			createCmd,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     serverAddressFlag,
				Aliases:  []string{"a"},
				Usage:    "Server address (can also be set via TASK_SERVER_ADDRESS env var)",
				Required: true,
				Sources:  cli.EnvVars("TASK_SERVER_ADDRESS"),
			},
			&cli.BoolFlag{
				Name:    insecureFlag,
				Aliases: []string{"plaintext"},
				Usage:   "Use plaintext connection (disable TLS). Can also be set via TASK_CLIENT_INSECURE env var.",
				Sources: cli.EnvVars("TASK_CLIENT_INSECURE"),
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
