package cmd

import (
	"dock/internal/container"

	"github.com/urfave/cli"
)

func init() {
	addSubCommands(runCmd)
}

var runCmd = cli.Command{
	Name:      "run",
	ShortName: "run",
	Usage:     "run docker container",
	UsageText: "run docker container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(ctx *cli.Context) error {
		// fork sub process, start sub process and quit
		initProc := container.NewParentProc(ctx.Bool("ti"), ctx.Args().Get(0))
		if err := initProc.Start(); err != nil {
			return err
		}

		// TODO: wait if daemon -d not set
		if err := initProc.Wait(); err != nil {
			return err
		}
		return nil
	},
}
