package main

import (
	"dock/internal/cmd"
	"log"
	"os"
)

func main() {

	app := cmd.App()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
