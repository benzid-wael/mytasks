package main

import (
	"fmt"
	"github.com/benzid-wael/mytasks/core/entities"
	"github.com/benzid-wael/mytasks/core/infrastructure"
)

func main() {
	var itemSequence entities.Sequence
	item := entities.NewNote(&itemSequence, "Golang is all about type", "", "@coding", "@golang")
	fmt.Println("item: ", item)
	item.Star()

	task := entities.NewTask(&itemSequence, "Read Accelerate book", "", "@reading")
	fmt.Printf("[%v] %v\n", task.Status, task.Title)

	task.TriggerEvent("start")
	fmt.Printf("[%v] %v\n", task.Status, task.Title)

	repo := infrastructure.NewItemRepository("~")
	repo.CreateTask(*task)
	repo.CreateNote(*item)

	items := repo.GetItems()
	fmt.Println("========== Storage ==========")
	for idx, item := range items {
		fmt.Printf("%v. %v\n", idx, item)
	}
}
