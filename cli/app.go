package cli

import (
	"fmt"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/benzid-wael/mytasks/tasks/usecases"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type MapFunction func(id int, c *cli.Context) error

func BulkFunc(ids []int, successMsgPrefix string, renderer Renderer, c *cli.Context, mapperFunc MapFunction) error {
	succeed := make([]string, 0, len(ids))
	failed := make([]string, 0, len(ids))
	for _, id := range ids {
		err := mapperFunc(id, c)
		if err != nil {
			fmt.Printf("Cannot process item: %v, error: %v\n", id, err)
			failed = append(failed, strconv.Itoa(id))
		} else {
			succeed = append(succeed, strconv.Itoa(id))
		}
	}

	fmt.Println()
	if len(succeed) > 0 {
		renderer.Success(successMsgPrefix + ": " + renderer.Colorify(strings.Join(succeed, ", "), GREY))
	}
	if len(failed) > 0 {
		renderer.Error("Failed Items: " + renderer.Colorify(strings.Join(failed, ", "), GREY))
	}

	return nil
}

func GetCliApp(config AppConfig) *cli.App {

	app := cli.App{
		Name:        "mytasks",
		Version:     "v0.0.0",
		Description: "📓 Manage your tasks and notes from the command line",
		Compiled:    time.Time{},
		Authors:     nil,
		Author:      "Wael Ben Zid El Guebsi",
		Email:       "benzid.wael@hotmail.fr",
	}

	itemRepository := infrastructure.NewItemRepository(config.DataDirectory)
	itemUseCase := usecases.NewItemUseCase(itemRepository)
	renderer := NewRenderer()
	itemPresenter := NewItemPresenter(renderer)

	app.Commands = []cli.Command{
		{
			Name:  "timeline",
			Usage: "Display Timeline View",
			Action: func(c *cli.Context) error {
				items := itemUseCase.GetItems()
				return itemPresenter.TimelineView(items)
			},
		},
		{
			Name:  "board",
			Usage: "Display Board View",
			Action: func(c *cli.Context) error {
				items := itemUseCase.GetItems()
				return itemPresenter.BoardView(items)
			},
		},
		{
			Name:    "note",
			Usage:   "Create note",
			Aliases: []string{"n"},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "title", Required: true},
				&cli.StringSliceFlag{Name: "tags"},
			},
			Action: func(c *cli.Context) error {
				tags := c.StringSlice("tags")
				title := c.String("title")
				note, err := itemUseCase.CreateNote(title, tags...)
				if err == nil {
					renderer.Success("Created note: " + renderer.Colorify(note.Id, GREY))
				} else {
					renderer.Error("Cannot create new note: " + err.Error())
				}
				return err
			},
		},
		{
			Name:    "task",
			Usage:   "Create task",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "title", Required: true},
				&cli.StringSliceFlag{Name: "tags"},
			},
			Action: func(c *cli.Context) error {
				tags := c.StringSlice("tags")
				title := c.String("title")
				note, err := itemUseCase.CreateTask(title, tags...)
				if err == nil {
					renderer.Success("Created task: " + renderer.Colorify(note.Id, GREY))
				} else {
					renderer.Error("Cannot create new task: " + err.Error())
				}
				return err
			},
		},
		{
			Name:  "clone",
			Usage: "Clone item",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				id := c.Int("id")
				item, err := itemUseCase.CloneItem(id)
				if err == nil {
					renderer.Success("Cloned item: " + renderer.Colorify(item.GetId(), GREY))
				} else {
					renderer.Error("Cannot clone item: " + err.Error())
				}
				return err
			},
		},
		{
			Name:    "edit",
			Usage:   "Edit item",
			Aliases: []string{"e"},
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
				&cli.StringFlag{Name: "title", Required: true},
				&cli.StringFlag{Name: "description"},
				&cli.StringSliceFlag{Name: "tags"},
			},
			Action: func(c *cli.Context) error {
				id := c.Int("id")
				title := c.String("title")
				description := c.String("description")
				tags := c.StringSlice("tags")
				err := itemUseCase.EditItem(id, title, description, nil, tags...)
				if err == nil {
					renderer.Success("Updated Item: " + renderer.Colorify(id, GREY))
				} else {
					renderer.Error("Cannot update item: " + err.Error())
				}
				return err
			},
		},
		{
			Name:  "star",
			Usage: "Star item",
			Flags: []cli.Flag{
				&cli.IntSliceFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				ids := c.IntSlice("id")

				starFunc := func(id int, c *cli.Context) error {
					starred := true
					return itemUseCase.EditItem(id, "", "", &starred, []string{}...)
				}
				return BulkFunc(ids, "Starred Items", renderer, c, starFunc)
			},
		},
		{
			Name:  "unstar",
			Usage: "Unstar item",
			Flags: []cli.Flag{
				&cli.IntSliceFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				ids := c.IntSlice("id")

				unstarFunc := func(id int, c *cli.Context) error {
					starred := false
					return itemUseCase.EditItem(id, "", "", &starred, []string{}...)
				}
				return BulkFunc(ids, "Unstarred Items", renderer, c, unstarFunc)
			},
		},
		{
			Name:    "archive",
			Usage:   "Archive item",
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				&cli.IntSliceFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				ids := c.IntSlice("id")

				archiveFunc := func(id int, c *cli.Context) error {
					return itemUseCase.ArchiveItem(id)
				}
				return BulkFunc(ids, "Archived Items", renderer, c, archiveFunc)
			},
		},
		{
			Name:    "restore",
			Usage:   "Restore item from archive",
			Aliases: []string{"r"},
			Flags: []cli.Flag{
				&cli.IntSliceFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				ids := c.IntSlice("id")

				restoreFunc := func(id int, c *cli.Context) error {
					return itemUseCase.RestoreItem(id)
				}
				return BulkFunc(ids, "Archived Items", renderer, c, restoreFunc)
			},
		},
		{
			Name:    "delete",
			Usage:   "Delete item",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				&cli.IntSliceFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				ids := c.IntSlice("id")

				deleteFunc := func(id int, c *cli.Context) error {
					return itemUseCase.DeleteItem(id)
				}
				return BulkFunc(ids, "Archived Items", renderer, c, deleteFunc)
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
