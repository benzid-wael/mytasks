package infrastructure

import (
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
)

type ItemRepositoryError struct {
	message string
}

type DoesNotExistError struct {
	message string
}

type AlreadyArchivedError struct {
	message string
}

func (error ItemRepositoryError) Error() string {
	return "ItemRepositoryError: " + error.message
}

func (error DoesNotExistError) Error() string {
	return "DoesNotExistError: " + error.message
}

func (error AlreadyArchivedError) Error() string {
	return "AlreadyArchivedError: " + error.message
}

type ItemRepository interface {
	StoreItem(id int, item interface{}) error
	CreateTask(task entities.Task) (entities.Task, error)
	CreateNote(note entities.Note) (entities.Note, error)

	GetItems() entities.ItemCollection
	GetItem(id int) entities.Manageable
	CloneItem(id int) (entities.Manageable, error)
	GetTaskById(id int) (*entities.Task, error)
	GetNoteById(id int) (*entities.Note, error)

	UpdateItem(id int, title *string, description *string, starred *bool, tags ...string) error
	DeleteItem(id int) error

	ArchiveItem(id int) error
	RestoreItem(id int) error
}
