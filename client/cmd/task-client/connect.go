package main

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
)

var connectCmd = &cli.Command{
	Name:  "connect",
	Usage: "Connect to a running Task Manager server.",
	Action: func(_ context.Context, c *cli.Command) error {
		args := c.Args().Slice()
		if len(args) != 1 {
			return errors.New("exactly one argument required with the URI of Task Manager server to connect to")
		}

		return nil
	},
}
