package usecases

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"sort"
)

type ItemUseCase interface {
	GetItems() []entities.Manageable
	CreateNote(title string, tags ...string) (*entities.Note, error)
	CreateTask(title string, tags ...string) (*entities.Task, error)
	EditItem(id int, title *string, description *string, tags ...string) error
	CopyItem(id int) (*entities.Note, error)
	ArchiveItem(id int) error
	RestoreItem(id int) error
	DeleteItem(id int) error
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

func (iuc *itemUseCase) EditItem(id int, title string, description string, starred *bool, tags ...string) error {
	newTitle := &title
	newDescription := &description
	if title == "" {
		newTitle = nil
	}
	if description == "" {
		newDescription = nil
	}
	return iuc.repository.UpdateItem(id, newTitle, newDescription, starred, tags...)
}

func (iuc *itemUseCase) GetItems() entities.ItemCollection {
	items := iuc.repository.GetItems()
	sort.Sort(items)
	return items
}

func (iuc *itemUseCase) CloneItem(id int) (entities.Manageable, error) {
	item, err := iuc.repository.CloneItem(id)
	return item, err
}

func (iuc *itemUseCase) ArchiveItem(id int) error {
	return iuc.repository.ArchiveItem(id)
}

func (iuc *itemUseCase) RestoreItem(id int) error {
	return iuc.repository.RestoreItem(id)
}

func (iuc *itemUseCase) DeleteItem(id int) error {
	return iuc.repository.ArchiveItem(id)
}
