package usecases

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
)

type ItemUseCase interface {
	CreateNote(title string, tags ...string) (*entities.Note, error)
	CreateTask(title string, tags ...string) (*entities.Task, error)
	EditItem(id int, title *string, description *string, tags ...string) error
}

type itemUseCase struct {
	repository infrastructure.ItemRepository
}

func NewItemUseCase(repository infrastructure.ItemRepository) *itemUseCase {
	return &itemUseCase{
		repository: repository,
	}
}

func (iuc *itemUseCase) CreateNote(title string, tags ...string) (*entities.Note, error) {
	item := entities.NewNote(title, "", tags...)
	note, err := iuc.repository.CreateNote(*item)
	return &note, err
}

func (iuc *itemUseCase) CreateTask(title string, tags ...string) (*entities.Task, error) {
	item := entities.NewTask(title, "", tags...)
	task, err := iuc.repository.CreateTask(*item)
	return &task, err
}

func (iuc *itemUseCase) EditItem(id int, title *string, description *string, tags ...string) error {
	return iuc.repository.UpdateItem(id, title, description, tags...)
}
