package main

import (
	"github.com/benzid-wael/mytasks/cli"
	notificator "github.com/benzid-wael/mytasks/notification"
)

var notify *notificator.Notificator

func main() {

	notify = notificator.New(notificator.Options{
		DefaultIcon: "icon/default.png",
		AppName:     "My test App",
	})

	notify.Push("title", "text", "/home/user/icon.png", notificator.UR_CRITICAL)
	cli.Main()
}
