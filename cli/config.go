package cli

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	DataDirectory         string `json:"data_directory"`
	DisplayCompletedTasks bool   `json:"display_completed_tasks"`
	DefaultDisplayMode    string `json:"default_display_mode"`
	DefaultWorkplace      string `json:"default_workplace"`
}

func GetAppConfig(path string, defaultDataDir string) AppConfig {
	appConfig := AppConfig{
		DataDirectory:         defaultDataDir,
		DisplayCompletedTasks: false,
		DefaultDisplayMode:    "board",
		DefaultWorkplace:      "",
	}
	path = tasks.ExpandPath(path)
	defer storeAppConfig(path, &appConfig) // nolint

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
