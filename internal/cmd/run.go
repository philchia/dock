package cmd

import (
	"dock/internal/container"
	"log"

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
		log.Println("run ", ctx.Args().Get(0), "tty:", ctx.Bool("ti"))
		// fork sub process, start sub process and quit
		initProc := container.NewParentProc(ctx.Bool("ti"), ctx.Args().Get(0))
		if err := initProc.Start(); err != nil {
			return err
		}

		if err := initProc.Wait(); err != nil {
			log.Println("wait:", err)
		}
		return nil
	},
}
