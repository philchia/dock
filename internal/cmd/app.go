package cmd

import (
	"github.com/urfave/cli"
)

var app = cli.NewApp()

// App get global app instance
func App() *cli.App {
	return app
}

func addSubCommands(cmds ...cli.Command) {
	app.Commands = append(app.Commands, cmds...)
}
