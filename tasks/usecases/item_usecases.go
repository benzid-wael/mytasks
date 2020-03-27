package usecases

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
)

type ItemUseCase interface {
	CreateNote(title string, tags ...string) (*entities.Note, error)
	CreateTask(title string, tags ...string) (*entities.Task, error)
	EditItem(id int, title *string, description *string, tags ...string) error
}

type itemUseCase struct {
	repository infrastructure.ItemRepository
	sequence   *value_objects.Sequence
}

func NewItemUseCase(repository infrastructure.ItemRepository, sequence *value_objects.Sequence) *itemUseCase {
	return &itemUseCase{
		repository: repository,
		sequence:   sequence,
	}
}

func (iuc *itemUseCase) CreateNote(title string, tags ...string) (*entities.Note, error) {
	note := entities.NewNote(iuc.sequence, title, "", tags...)
	err := iuc.repository.CreateNote(*note)
	return note, err
}

func (iuc *itemUseCase) CreateTask(title string, tags ...string) (*entities.Task, error) {
	task := entities.NewTask(iuc.sequence, title, "", tags...)
	err := iuc.repository.CreateTask(*task)
	return task, err
}

func (iuc *itemUseCase) EditItem(id int, title *string, description *string, tags ...string) error {
	return iuc.repository.UpdateItem(id, title, description, tags...)
}
