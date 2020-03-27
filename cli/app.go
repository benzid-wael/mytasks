package cli

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/benzid-wael/mytasks/tasks/usecases"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

type AppConfig struct {
	DataDirectory         string `json:"data_directory"`
	DisplayCompletedTasks bool   `json:"display_completed_tasks"`
	DefaultDisplayMode    string `json:"default_display_mode"`
}

type AppState struct {
	ItemSequence value_objects.Sequence
}

func GetAppConfig(path string, defaultDataDir string) AppConfig {
	appConfig := AppConfig{
		DataDirectory:         defaultDataDir,
		DisplayCompletedTasks: false,
		DefaultDisplayMode:    "board",
	}
	path = tasks.ExpandPath(path)
	defer storeAppConfig(path, &appConfig)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Cannot load config file: ", err.Error())
	} else {
		err = json.Unmarshal(data, &appConfig)
		if err != nil {
			log.Println("Invalid config file: ", err.Error())
		}
	}

	return appConfig
}

func storeAppConfig(path string, config *AppConfig) error {
	payload, err := json.Marshal(*config)
	if err != nil {
		fmt.Println("Cannot marshall items: ", err)
		return err
	}

	return ioutil.WriteFile(path, payload, 0644)
}

func GetAppState(path string) AppState {
	appState := AppState{ItemSequence: *value_objects.NewSequence(0)}
	path = tasks.ExpandPath(path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Cannot load mytasks.state: ", err.Error())
	} else {
		err = json.Unmarshal(data, &appState)
		if err != nil {
			log.Println("Invalid mytasks.state file: ", err.Error())
		}
	}
	return appState
}

func storeAppState(path string, state *AppState) error {
	payload, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Cannot marshall AppState: %v, err: %v\n", state, err)
		return err
	}
	return ioutil.WriteFile(path, payload, 0644)
}

func GetCliApp(config AppConfig, state *AppState) *cli.App {

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
	itemUseCase := usecases.NewItemUseCase(itemRepository, &state.ItemSequence)
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
	dataDir := tasks.ExpandPath("~/.mytasks")
	appConfigPath := tasks.ExpandPath("~/.mytasks.json")
	appStatePath := path.Join(dataDir, ".mytasks.state")
	appConfig := GetAppConfig(appConfigPath, dataDir)
	appState := GetAppState(appStatePath)
	defer storeAppState(appStatePath, &appState)
	app := GetCliApp(appConfig, &appState)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Cannot run command: %v. Original error: %v", os.Args, err)
	}
}
