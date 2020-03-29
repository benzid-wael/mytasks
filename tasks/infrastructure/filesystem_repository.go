package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"github.com/benzid-wael/mytasks/tasks/domain/value_objects"
	"io/ioutil"
	"os"
	path2 "path"
	"strconv"
	"time"
)

type FilesystemItemRepository struct {
	DataDirectory    string `json:"data_directory"`
	StorageDir       string `json:"repository_dir"`
	ArchiveDir       string `json:"archive_dir"`
	itemSeqStatePath string
}

func NewItemRepository(dataDir string) *FilesystemItemRepository {
	dataDir = tasks.ExpandPath(dataDir)
	StortageDir := path2.Join(dataDir, "./repository")
	ArchiveDir := path2.Join(dataDir, "./archive")

	tasks.CreateDirIfNotExist(StortageDir)
	tasks.CreateDirIfNotExist(ArchiveDir)

	return &FilesystemItemRepository{
		DataDirectory:    dataDir,
		StorageDir:       StortageDir,
		ArchiveDir:       ArchiveDir,
		itemSeqStatePath: path2.Join(dataDir, ".item.sequence.state"),
	}
}

func getKey(id int) string {
	return strconv.Itoa(id)
}

func loadItems(dir string) map[string]map[string]interface{} {
	path := path2.Join(dir, "items.json")
	items := make(map[string]map[string]interface{})

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("[1] Cannot open: %v, err: %v\n", path, err)
		return items
	} else if err != nil {
		fmt.Printf("[2] Cannot open: %v, err: %v\n", path, err)
		panic(err)
	}
	if data, err := ioutil.ReadFile(path); err != nil {
		fmt.Printf("[3] Cannot load items from directory: %v, err: %v\n ", dir, err)
		panic(err)
	} else {
		if err := json.Unmarshal(data, &items); err != nil {
			fmt.Printf("[4] Cannot load items in directory: %v, err: %v\n", path, err)
			panic(err)
		}
		return items
	}
}

func storeItems(items map[string]map[string]interface{}, dir string) error {
	if data, err := json.Marshal(items); err != nil {
		fmt.Println("Cannot marshall items: ", err)
		return err
	} else {
		path := path2.Join(dir, "items.json")
		err := ioutil.WriteFile(path, data, 0644)
		if err != nil {
			fmt.Printf("Cannot store items in directory: %v, err: %v\n", path, err)
			return ItemRepositoryError{"Cannot store items in directory: " + dir + ". Original error: " + err.Error()}
		}
	}
	return nil
}

func (repository *FilesystemItemRepository) store(id string, item interface{}, dir string) error {
	items := loadItems(dir)

	if data, err := tasks.ToMap(item); err == nil {
		items[id] = *data
	}

	return storeItems(items, dir)
}

func (repository *FilesystemItemRepository) storeItem(id string, item interface{}) error {
	return repository.store(id, item, repository.StorageDir)
}

func (repository *FilesystemItemRepository) archiveItem(id string, item interface{}) error {
	return repository.store(id, item, repository.ArchiveDir)
}

func (repository *FilesystemItemRepository) GetNextId() int {
	sequence := value_objects.NewSequenceFromFS(repository.itemSeqStatePath)
	if sequence == nil {
		sequence = value_objects.NewSequence(0)
	}
	defer value_objects.SaveSequence(repository.itemSeqStatePath, sequence)
	sequence.Next()
	return sequence.Current()
}

func (repository *FilesystemItemRepository) CreateTask(task entities.Task) (entities.Task, error) {
	task.Id = repository.GetNextId()
	err := repository.storeItem(getKey(task.Id), task)
	return task, err
}

func (repository *FilesystemItemRepository) CreateNote(note entities.Note) (entities.Note, error) {
	note.Id = repository.GetNextId()
	err := repository.storeItem(getKey(note.Id), note)
	return note, err
}

func (repository *FilesystemItemRepository) GetItems() entities.ItemCollection {
	items := loadItems(repository.StorageDir)
	result := make(entities.ItemCollection, len(items))

	index := 0
	for _, item := range items {
		if item != nil {
			result[index] = entities.CreateItem(item)
			index++
		}
	}

	return result
}

func (repository *FilesystemItemRepository) getItem(items map[string]map[string]interface{}, id int) (map[string]interface{}, error) {
	key := getKey(id)
	if item, ok := items[key]; ok {
		return item, nil
	}
	return nil, DoesNotExistError{"Cannot find any item with ID: " + key}
}

func (repository *FilesystemItemRepository) GetItem(id int) *entities.Manageable {
	items := loadItems(repository.StorageDir)
	data, err := repository.getItem(items, id)
	if err != nil {
		panic(err)
	}
	item := entities.CreateItem(data)
	return &item
}

func (repository *FilesystemItemRepository) CloneItem(id int) (entities.Manageable, error) {
	items := loadItems(repository.StorageDir)
	data, err := repository.getItem(items, id)
	if err != nil {
		return nil, err
	}

	data["created_at"] = time.Now().Format(time.RFC3339)
	newId := repository.GetNextId()
	data["id"] = float64(newId)
	err = repository.storeItem(getKey(newId), data)
	if err != nil {
		return nil, err
	}

	item := entities.CreateItem(data)
	return item, nil
}

func (repository *FilesystemItemRepository) UpdateItem(id int, title *string, description *string, starred *bool, tags ...string) error {
	items := loadItems(repository.StorageDir)

	if item, err := repository.getItem(items, id); err == nil {
		if title != nil {
			item["title"] = *title
		}
		if description != nil {
			item["description"] = *description
		}
		if len(tags) > 0 {
			item["tags"] = tags
		}
		if starred != nil {
			item["is_starred"] = starred
		}

		// Update collection
		key := getKey(id)
		items[key] = item

		// store new version
		return storeItems(items, repository.StorageDir)
	} else {
		return err
	}
}

func (repository *FilesystemItemRepository) DeleteItem(id int) error {
	items := loadItems(repository.StorageDir)
	if _, err := repository.getItem(items, id); err != nil {
		return err
	}
	delete(items, getKey(id))

	// store new version
	return storeItems(items, repository.StorageDir)
}

func (repository *FilesystemItemRepository) ArchiveItem(id int) error {
	items := loadItems(repository.StorageDir)
	item, err := repository.getItem(items, id)
	if err != nil {
		return err
	}
	archivedItems := loadItems(repository.ArchiveDir)
	_, err2 := repository.getItem(archivedItems, id)
	key := getKey(id)
	if err2 == nil {
		return AlreadyArchivedError{
			message: "There is already ab archived item with the same ID: " + key,
		}
	}

	delete(items, key)
	err3 := storeItems(items, repository.StorageDir)
	if err3 != nil {
		return ItemRepositoryError{"Cannot store items in the storage directory. Original error: " + err3.Error()}
	}

	return repository.archiveItem(key, item)
}

func (repository *FilesystemItemRepository) RestoreItem(id int) error {
	archivedItems := loadItems(repository.ArchiveDir)
	item, err := repository.getItem(archivedItems, id)
	if err != nil {
		return err
	}

	key := getKey(id)
	delete(archivedItems, key)
	err3 := storeItems(archivedItems, repository.StorageDir)
	if err3 != nil {
		return ItemRepositoryError{"Cannot store items in the archive directory. Original error: " + err3.Error()}
	}

	return repository.storeItem(key, item)
}
