package main

import (
	"github.com/benzid-wael/mytasks/cli"
	"github.com/benzid-wael/mytasks/tasks"
	"os"
)

func main() {
	dataDir := "~/.mytasks"
	appConfigPath := tasks.ExpandPath("~/.mytasks.json")
	appConfig := cli.GetAppConfig(appConfigPath, dataDir)
	app := cli.GetCliApp(appConfig)

	app.Run(os.Args)
}
