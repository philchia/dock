package cmd

import (
	"dock/internal/container"
	"dock/internal/subsystem"
	"strings"

	"github.com/google/uuid"

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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cupshare limit",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "root dir",
			Value: "/",
		},
	},
	Action: func(ctx *cli.Context) error {
		// fork sub process, start sub process and quit
		initProc, w, err := container.NewParentProc(ctx.Bool("ti"), ctx.String("root"))
		if err != nil {
			return err
		}

		if err := initProc.Start(); err != nil {
			return err
		}

		// write container cmd to sub process
		if _, err := w.WriteString(strings.Join(ctx.Args(), " ")); err != nil {
			return err
		}
		w.Close()

		conf := &subsystem.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CPUShare:    ctx.String("cpushare"),
		}

		cgroupManager := subsystem.NewCgroupManager(uuid.New().String())
		defer cgroupManager.Destroy()

		if err := cgroupManager.Set(conf); err != nil {
			return err
		}

		if err := cgroupManager.Apply(initProc.Process.Pid); err != nil {
			return err
		}

		// TODO: wait if daemon -d not set
		if err := initProc.Wait(); err != nil {
			return err
		}
		return nil
	},
}
