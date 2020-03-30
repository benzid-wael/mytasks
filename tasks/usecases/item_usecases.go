package usecases

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/infrastructure"
	"sort"
)

type ItemUseCase interface {
	CreateNote(title string, tags ...string) (*entities.Note, error)
	CreateTask(title string, tags ...string) (*entities.Task, error)
	GetItems() entities.ItemCollection
	GetItem(id int) entities.Manageable
	GetNoteById(id int) (*entities.Note, error)
	GetTaskById(id int) (*entities.Task, error)
	EditItem(id int, title string, description string, starred *bool, tags ...string) error
	CloneItem(id int) (entities.Manageable, error)
	ArchiveItem(id int) error
	RestoreItem(id int) error
	DeleteItem(id int) error
	TriggerEvent(id int, event string) error
}

type itemUseCase struct {
	repository infrastructure.ItemRepository
}

func NewItemUseCase(repository infrastructure.ItemRepository) ItemUseCase {
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
	item := entities.NewTask(title, "", entities.ToDo, tags...)
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

func (iuc *itemUseCase) GetItem(id int) entities.Manageable {
	item := iuc.repository.GetItem(id)
	return item
}

func (iuc *itemUseCase) GetItems() entities.ItemCollection {
	items := iuc.repository.GetItems()
	sort.Sort(items)
	return items
}

func (iuc *itemUseCase) GetNoteById(id int) (*entities.Note, error) {
	return iuc.repository.GetNoteById(id)
}

func (iuc *itemUseCase) GetTaskById(id int) (*entities.Task, error) {
	return iuc.repository.GetTaskById(id)
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

func (iuc *itemUseCase) TriggerEvent(id int, event string) error {
	item, err := iuc.GetTaskById(id)
	if err != nil {
		return err
	}

	err = item.TriggerEvent(event)
	if err != nil {
		return err
	}

	return iuc.repository.StoreItem(id, item)
}
