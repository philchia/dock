package cmd

import (
	"dock/internal/container"

	"github.com/urfave/cli"
)

func init() {
	addSubCommands(initCmd)
}

var initCmd = cli.Command{
	Name:      "init",
	ShortName: "init",
	Usage:     "do not call directly",
	UsageText: "do not call directly",
	Hidden:    true,

	Action: func(ctx *cli.Context) error {
		// run cmd as init process in container
		container.RunContainerInitProc()
		return nil
	},
}
