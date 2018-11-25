package cmd

import (
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

	Action: func(ctx *cli.Context) error {
		log.Println("run command")
		return nil
	},
}
