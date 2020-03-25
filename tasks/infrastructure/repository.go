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
	CreateTask(task entities.Task) error
	CreateNote(note entities.Note) error

	GetItem(id int) *entities.Manageable
	GetItems() []entities.Manageable

	UpdateItem(id int, title *string, description *string, tags ...string) error
	DeleteItem(id int) error

	ArchiveItem(id int) error
	RestoreItem(id int) error
}
