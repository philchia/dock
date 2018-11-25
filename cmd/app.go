package cmd

import (
	"github.com/urfave/cli"
)

var app = cli.NewApp()

// App get global app instance
func App() *cli.App {
	return app
}

func addFlags(flag cli.Flag) {
	app.Flags = append(app.Flags, flag)
}

func addSubCommands(cmds ...cli.Command) {
	app.Commands = append(app.Commands, cmds...)
}
