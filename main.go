package main

import (
	"dock/cmd"
	"log"
	"os"
)

func main() {

	app := cmd.App()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
