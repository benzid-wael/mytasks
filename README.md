# mytasks
 MyTasks is a simple solution to manage your tasks and notes across multiple boards from within your terminal.

## Features
 - Organize your tasks and notes
 - Lightweight tool
 - Manage tasks effectively
 - Archive notes / tasks
 - Restore archive notes / tasks
 - Favorite mechanism
 - Support of multiple display modes
 - Simple usage syntax

## Roadmap
 - Search & filter mechanism
 - Tasks due date & custom notification
 - Custom template
 - Profile support
 
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

### Display available tasks and notes

As per now, we support two display modes:

* Timeline view where items are grouped per creation date
* Board view where items are grouped per board. A board is identified by a unique tag