# mytasks
MyTasks is a simple solution to manage your tasks and notes across multiple boards from within your terminal.
MyTasks is heavily inspired by [taskbook](https://github.com/klaussinani/taskbook) with some minor improvements.

## Features
 - Organize your tasks and notes
 - Lightweight tool
 - Manage tasks effectively
 - Archive notes / tasks
 - Restore archive notes / tasks
 - Favorite mechanism
 - Support of multiple display modes
 - Simple usage syntax
 - Workspace support
 - Search & filter mechanism
 - Tasks due date & custom notification

## Roadmap
 - Custom template

## Installation

 You can install the package using go tool

 ```shell script
$ go get -u github.com/benzid-wael/mytasks
```

IF you want to install only the cli tool

 ```shell script
$ go get -u github.com/benzid-wael/mytasks/cmd/mytasks
```

I recommend installing it using [bingo](https://github.com/TekWizely/bingo)

 ```shell script
$ bingo install github.com/benzid-wael/mytasks/cmd/mytasks@v0.0.1
```

Later you can uninstall it using this command

 ```shell script
$ bingo uninstall mytasks
```

## Usage

### Creating new task

 ```shell script
$ mytasks task --title "Learn golang" --tags coding --tags golang
```


### Creating new note

 ```shell script
$ mytasks note --title "Golang is all about types" --tags golang
```

### Special Boards

When you call mytasks without specifying any command, mytasks will display all tasks in any of the following board
and in the following order:
- `Today List` contains list of most important tasks (high priority ones)
- `Tomorrow List` contains list of less important tasks, that you should pick from when you don't have anything in your
Today List
- `Next List`: Contains list of anything else (low-priority tasks)

Please note that in this version, this feature is dump and we don't generate this list dynamically. Instead, you need
to manage it yourself.

To create a task in your today's list:
 ```shell script
$ mytasks note --title "Golang is all about types" --tags Today
```

### Display available tasks and notes

As per now, we support two display modes:

* Timeline view where items are grouped per creation date
* Board view where items are grouped per board. A board is identified by a unique tag

To display the timeline view, run this command
 ```shell script
$ mytasks timeline
```

In order to display the board view, you need to use `board` command
 ```shell script
$ mytasks board
```

For further control, you need to use `list` command which allow you to filter items by specific attributes.
As per current version, we support the following attributes:
- `status`: Keep only items with given status
- `type`: Keeps only items with given type. Supported types: `note` and `task`
- `due-date`: Keeps only items due on the given due date, format: `YYYY-MM-DD`
- `created-before`: Keeps items created before given date, format: `YYYY-MM-DD`
- `created-after`: Keeps items created before given date, format: `YYYY-MM-DD`
- `exclude`: Exclude given boards
- `tags`: Include only given boards

### Install Notification Jobs
MyTasks comes with a command to remind you with tasks about to be due

To get notifcations about tasks to be due within the next 2 hours, you need to run the following command
```shell script
$ mytasks-notifier notify --before 120
```

In practice, you need to add cron jobs using your cron daemon or any similar system. To add job that remind me every
15 minutes about task due within 2 hours, you need to do the following:

 ```shell script
$ crontab -e
```

and then add the following line and save:

```
*/15 * * * * mytasks-notifier notify --before 120
```

For further information about crontab expression, you can check [crontab.guru](https://crontab.guru/) website
