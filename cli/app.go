package cli

import (
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/benzid-wael/mytasks/tasks/usecases"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

func GetCliApp(config AppConfig) *cli.App {

	app := cli.App{
		Name:        "mytasks",
		Version:     "v0.0.0",
		Description: "Manage your tasks and notes from the command line",
		Compiled:    time.Time{},
		Authors:     nil,
		Author:      "Wael Ben Zid El Guebsi",
		Email:       "benzid.wael@hotmail.fr",
	}

	itemRepository := infrastructure.NewItemRepository(config.DataDirectory)
	itemUseCase := usecases.NewItemUseCase(itemRepository)
	renderer := NewRenderer()

	app.Commands = []cli.Command{
		{
			Name:        "note",
			Aliases:     []string{"n"},
			Description: "Create note",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{Name: "tags"},
				&cli.StringSliceFlag{Name: "title", Required: true},
			},
			Action: func(c *cli.Context) error {
				tags := c.StringSlice("tags")
				title := c.String("title")
				note, err := itemUseCase.CreateNote(title, tags...)
				if err == nil {
					renderer.Success("Created note: " + renderer.Colorify(note.Id, GREY))
				}
				return err
			},
		},
	}

	return &app
}

func Main() {
	dataDir := "~/.mytasks"
	appConfigPath := tasks.ExpandPath("~/.mytasks.json")
	appConfig := GetAppConfig(appConfigPath, dataDir)
	app := GetCliApp(appConfig)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Cannot run command: %v. Original error: %v", os.Args, err)
	}
}
