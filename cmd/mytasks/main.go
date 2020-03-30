package main

import (
	"github.com/benzid-wael/mytasks/cli"
	"log"
	"os"
)

func main() {
	app := cli.GetCliApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
