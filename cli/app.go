package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/benzid-wael/mytasks/tasks/usecases"
	"github.com/urfave/cli/v2"
	"path"
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
			failed = append(failed, strconv.Itoa(id)) // nolint
		} else {
			succeed = append(succeed, strconv.Itoa(id)) // nolint
		}
	}

	fmt.Println()
	if len(succeed) > 0 {
		return renderer.Success(successMsgPrefix + ": " + renderer.Colorify(strings.Join(succeed, ", "), GREY))
	}
	if len(failed) > 0 {
		return renderer.Error("Failed Items: " + renderer.Colorify(strings.Join(failed, ", "), GREY))
	}
	return nil
}

func GetItemUseCase(workspace string, config AppConfig) usecases.ItemUseCase {
	dataDir := config.DataDirectory
	if workspace != "" {
		dataDir = path.Join(dataDir, workspace)
	}
	itemRepository := infrastructure.NewItemRepository(dataDir)
	return usecases.NewItemUseCase(itemRepository)
}

func GetCliApp() *cli.App {
	dataDir := "~/.mytasks"
	appConfigPath := tasks.ExpandPath("~/.mytasks.json")
	config := GetAppConfig(appConfigPath, dataDir)

	renderer := NewRenderer()
	itemPresenter := NewItemPresenter(renderer)

	app := cli.App{
		Name:                 "mytasks",
		Version:              "v0.0.1",
		EnableBashCompletion: true,
		Description:          "📓 Manage your tasks and notes from the command line",
		Compiled:             time.Time{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "workspace",
				Usage:   "Workspace name",
				Value:   config.DefaultWorkplace,
				EnvVars: []string{"MYTASKS_WORKSPACE"},
			},
		},
		Action: func(c *cli.Context) error {
			itemUseCase := GetItemUseCase(c.String("workspace"), config)
			tags := []string{"Today", "Tomorrow", "Next List"}
			items := itemUseCase.GetItems().FilterByTags(tags...).FilterPending()
			return itemPresenter.BoardView(items, tags...)
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "timeline",
			Usage: "Display Timeline View",
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				items := itemUseCase.GetItems()
				return itemPresenter.TimelineView(items)
			},
		},
		{
			Name:  "board",
			Usage: "Display Board View",
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				tags := c.StringSlice("tags")
				title := c.String("title")
				note, err := itemUseCase.CreateNote(title, tags...)
				if err == nil {
					renderer.Success("Created note: " + renderer.Colorify(note.Id, GREY)) // nolint
				} else {
					return errors.New("Cannot create new note: " + err.Error()) // nolint
				}
				return nil
			},
		},
		{
			Name:    "task",
			Usage:   "Create task",
			Aliases: []string{"t"},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "title", Required: true},
				&cli.StringFlag{Name: "priority"},
				&cli.StringFlag{Name: "due-date"},
				&cli.StringSliceFlag{Name: "tags"},
			},
			Action: func(c *cli.Context) error {
				data := make(map[string]interface{}, 1)
				if c.String("priority") != "" {
					priority, err := getPriority(c.String("priority"))
					if err != nil {
						return err
					}
					data["priority"] = priority
				}

				if c.String("due-date") != "" {
					value, err := getDate(c.String("due-date"))
					if err != nil {
						return err
					}
					data["due_date"] = value
				}

				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				tags := c.StringSlice("tags")
				title := c.String("title")
				task, err2 := itemUseCase.CreateTask(title, tags...)
				if err2 == nil {
					if len(data) > 0 {
						err2 = itemUseCase.EditItem(task.Id, data)
					}
					renderer.Success("Created task: " + renderer.Colorify(task.Id, GREY)) // nolint
				} else {
					return errors.New("Cannot create new task: " + err2.Error()) // nolint
				}
				return err2
			},
		},
		{
			Name:  "clone",
			Usage: "Clone item",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				item, err := itemUseCase.CloneItem(id)
				if err == nil {
					renderer.Success("Cloned item: " + renderer.Colorify(item.GetId(), GREY)) // nolint
				} else {
					return errors.New("Cannot clone item: " + err.Error()) // nolint
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
				&cli.StringFlag{Name: "title"},
				&cli.StringFlag{Name: "description"},
				&cli.StringFlag{Name: "priority"},
				&cli.StringFlag{Name: "due-date"},
				&cli.StringSliceFlag{Name: "tags"},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				data := map[string]interface{}{
					"title":       c.String("title"),
					"description": c.String("description"),
					"tags":        c.StringSlice("tags"),
				}
				if c.String("priority") != "" || c.String("due-date") != "" {
					if c.String("priority") != "" {
						priority, err := getPriority(c.String("priority"))
						if err != nil {
							return err
						}
						data["priority"] = priority
					}
					if c.String("due-date") != "" {
						dueDate, err := getDate(c.String("due-date"))
						if err != nil {
							return err
						}
						data["due_date"] = dueDate
					}
					_, err := itemUseCase.GetTaskById(id)
					if err != nil {
						return err
					}
				}
				err := itemUseCase.EditItem(id, data)
				if err == nil {
					renderer.Success("Updated Item: " + renderer.Colorify(id, GREY)) // nolint
				} else {
					return errors.New("Cannot update item: " + err.Error()) // nolint
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				ids := c.IntSlice("id")
				data := map[string]interface{}{
					"is_starred": true,
				}

				starFunc := func(id int, c *cli.Context) error {
					return itemUseCase.EditItem(id, data)
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				ids := c.IntSlice("id")
				data := map[string]interface{}{
					"is_starred": false,
				}

				unstarFunc := func(id int, c *cli.Context) error {
					return itemUseCase.EditItem(id, data)
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
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
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				ids := c.IntSlice("id")

				deleteFunc := func(id int, c *cli.Context) error {
					return itemUseCase.DeleteItem(id)
				}
				return BulkFunc(ids, "Archived Items", renderer, c, deleteFunc)
			},
		},
		{
			Name:  "copy",
			Usage: "Copy details to clipboard",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
				&cli.BoolFlag{Name: "title"},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				item := itemUseCase.GetItem(id)
				value := ""
				if c.Bool("title") {
					value = item.GetTitle()
				} else {
					bytes, _ := json.Marshal(item)
					value = string(bytes)
				}
				err := clipboard.WriteAll(value)
				if err != nil {
					renderer.Success("Item's information copied successfully") // nolint
				} else {
					return errors.New("cannot copy item's information")
				}
				return err
			},
		},
		{
			Name:  "start",
			Usage: "Start task",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				err := itemUseCase.TriggerEvent(id, "start")
				if err == nil {
					renderer.Success("Task started") // nolint
				}
				return err
			},
		},
		{
			Name:  "pause",
			Usage: "Pause task",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				err := itemUseCase.TriggerEvent(id, "stop")
				if err == nil {
					renderer.Success("Task paused") // nolint
				}
				return err
			},
		},
		{
			Name:  "complete",
			Usage: "Complete task",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				err := itemUseCase.TriggerEvent(id, "complete")
				if err == nil {
					renderer.Success("Task marked as completed") // nolint
				}
				return err
			},
		},
		{
			Name:  "cancel",
			Usage: "Cancel task",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "id", Required: true},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				id := c.Int("id")
				err := itemUseCase.TriggerEvent(id, "cancel")
				if err == nil {
					renderer.Success("Task cancelled") // nolint
				}
				return err
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list items matching given criteria",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "status"},
				&cli.StringFlag{Name: "type"},
				&cli.StringSliceFlag{Name: "tags"},
				&cli.StringFlag{Name: "view", Usage: "Display mode: timeline or board", Value: "timeline"},
				&cli.StringFlag{Name: "due-date"},
				&cli.StringFlag{Name: "created-before"},
				&cli.StringFlag{Name: "created-after"},
				&cli.StringSliceFlag{Name: "exclude", Usage: "Tags to be excluded"},
			},
			Action: func(c *cli.Context) error {
				view := c.String("view")
				if view != "timeline" && view != "board" {
					return errors.New("Unsupported view: " + view)
				}
				itemUseCase := GetItemUseCase(c.String("workspace"), config)
				items := itemUseCase.GetItems()

				status := c.String("status")
				if status == "pending" {
					items = items.FilterPending()
				} else if status != "" {
					items = items.FilterByStatus(status)
				}

				kind := c.String("type")
				if kind != "" {
					items = items.FilterByType(kind)
				}

				tags := c.StringSlice("tags")
				if len(tags) > 0 {
					items = items.FilterByTags(tags...)
				}

				dueDateString := c.String("due-date")
				if dueDateString != "" {
					dueDate, err := getDate(dueDateString)
					if err != nil {
						return err
					}
					items = items.Filter(func(item entities.Manageable) bool {
						itemDueDate := item.GetDueDate()
						return itemDueDate != nil && dueDate.Equal(*itemDueDate)
					})
				}

				items = items.FilterByCreationDate(
					getDateOrNil(c.String("created-after")),
					getDateOrNil(c.String("created-before")),
				)

				excludeTags := c.StringSlice("exclude")
				for _, tag := range excludeTags {
					hasTagFilter := func(item entities.Manageable) bool {
						return item.HasTag(tag)
					}

					items = items.Exclude(hasTagFilter)
				}

				if view == "timeline" {
					return itemPresenter.TimelineView(items)
				} else if view == "board" {
					return itemPresenter.BoardView(items)
				}
				return nil
			},
		},
	}

	return &app
}
