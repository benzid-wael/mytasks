package main

import (
	"fmt"
	"github.com/benzid-wael/mytasks/core/entities"
)

func main() {
	var itemSequence entities.Sequence
	item := entities.NewNote(&itemSequence, "Golang is all about type", "@coding", "@golang")

	fmt.Println("item: ", item)
	item.Star()
}
