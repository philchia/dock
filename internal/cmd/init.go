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
		cmd := ctx.Args().Get(0)
		// run cmd as init process in container
		container.RunContainerInitProc(cmd)
		return nil
	},
}
