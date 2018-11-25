package cmd

import (
	"github.com/urfave/cli"
)

var app = &cli.App{
	Name:                 "dock",
	Version:              "1.0",
	EnableBashCompletion: true,
}

// App get global app instance
func App() *cli.App {
	return app
}

func addSubCommands(cmds ...cli.Command) {
	app.Commands = append(app.Commands, cmds...)
}
