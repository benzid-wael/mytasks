package main

import (
	"fmt"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"github.com/benzid-wael/mytasks/tasks/usecases"
)

func main() {
	var itemSequence value_objects.Sequence
	item := entities.NewNote(&itemSequence, "Golang is all about type", "", "@coding", "@golang")
	fmt.Println("item: ", item)
	item.Star()

	task := entities.NewTask(&itemSequence, "Read Accelerate book", "", "@reading")
	fmt.Printf("[%v] %v\n", task.Status, task.Title)

	task.TriggerEvent("start")
	fmt.Printf("[%v] %v\n", task.Status, task.Title)

	fmt.Println("========== Repository ==========")
	repo := infrastructure.NewItemRepository("~")
	repo.CreateTask(*task)
	repo.CreateNote(*item)

	fmt.Println("========== Usecase ==========")
	usecase := usecases.NewItemUseCase(repo, itemSequence)
	note2, _ := usecase.CreateNote("You should always write good code. You should always use abstraction", "@cleanCode")
	fmt.Printf("%v\n", note2.Title)

	items := repo.GetItems()
	fmt.Println("========== Storage ==========")
	for idx, item := range items {
		fmt.Printf("%v. %v\n", idx, item)
	}
}
