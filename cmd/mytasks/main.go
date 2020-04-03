package main

import (
	"github.com/benzid-wael/mytasks/cli"
	"os"
)

func main() {
	renderer := cli.NewRenderer()
	app := cli.GetCliApp()

	err := app.Run(os.Args)
	if err != nil {
		renderer.Error(err.Error()) // nolint
	}
}
