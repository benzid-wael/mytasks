package main

import (
	"fmt"
	cli2 "github.com/benzid-wael/mytasks/cli"
	notificator "github.com/benzid-wael/mytasks/notification"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {

	dataDir := "~/.mytasks"
	appConfigPath := tasks.ExpandPath("~/.mytasks.json")
	config := cli2.GetAppConfig(appConfigPath, dataDir)

	app := cli.App{
		Name:        "mytasks-notifier",
		Version:     "v0.0.1",
		Description: "Be reminded before your task is due!",
		Compiled:    time.Time{},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "workspace",
				Usage:  "Workspace name",
				Value:  config.DefaultWorkplace,
				EnvVar: "MYTASKS_WORKSPACE",
			},
		},
	}

	notify := notificator.New(notificator.Options{
		DefaultIcon: "/opt/mytasks/logo.png",
		AppName:     "MyTasks",
	})

	app.Commands = []cli.Command{
		{
			Name:  "notify",
			Usage: "Create notifications for your tasks",
			Flags: []cli.Flag{
				&cli.IntFlag{Name: "before", Usage: "How many minutes to be reminded before the task is due"},
			},
			Action: func(c *cli.Context) error {
				itemUseCase := cli2.GetItemUseCase(c.GlobalString("workspace"), config)
				items := itemUseCase.GetItems()

				now := time.Now()
				later := now.Add(time.Hour * time.Duration(24*c.Int("before")+1))
				for _, item := range items {
					const dateLayout string = "January 2, 2006"
					dueDate := item.GetDueDate()
					if dueDate == nil {
						fmt.Printf("Task %v does not have a due date\n", item.GetId())
						continue
					}

					status := item.GetStatus()
					if status == string(entities.Completed) || status == string(entities.Cancelled) {
						fmt.Printf("Task %v is in a final state\n", item.GetId())
						continue
					}

					var title, description string
					if now.After(*dueDate) {
						title = fmt.Sprintf("Ⓘ Task %v is overdue", item.GetId())
						description = fmt.Sprintf("%v - %v", item.GetTitle(), status)
					} else if later.After(*dueDate) {
						title = fmt.Sprintf("ℹ Task %v is due on %v", item.GetId(), dueDate.Format(dateLayout))
						description = fmt.Sprintf("%v - %v", item.GetTitle(), status)
					} else {
						continue
					}

					notify.Push(title, description, "", notificator.UR_NORMAL) // nolint
				}

				return nil
			},
		},
	}

	app.Run(os.Args) // nolint
}
