package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/benzid-wael/mytasks/tasks"
	"github.com/benzid-wael/mytasks/tasks/domain/entities"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type FilesystemItemRepository struct {
	MainAppDir string `json:"main_app_dir"`
	StorageDir string `json:"repository_dir"`
	ArchiveDir string `json:"archive_dir"`
}

func NewItemRepository(appDir string) *FilesystemItemRepository {

	DataDir := path.Join(appDir, "./mytasks")
	StortageDir := path.Join(DataDir, "./repository")
	ArchiveDir := path.Join(DataDir, "./archive")

	tasks.CreateDirIfNotExist(StortageDir)
	tasks.CreateDirIfNotExist(ArchiveDir)

	return &FilesystemItemRepository{
		MainAppDir: appDir,
		StorageDir: StortageDir,
		ArchiveDir: ArchiveDir,
	}
}

func getKey(id int) string {
	return strconv.Itoa(id)
}

func loadItems(dir string) map[string]map[string]interface{} {
	path := path.Join(dir, "items.json")
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
		path := path.Join(dir, "items.json")
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

func (repository *FilesystemItemRepository) CreateTask(task entities.Task) error {
	return repository.storeItem(getKey(task.Id), task)
}

func (repository *FilesystemItemRepository) CreateNote(note entities.Note) error {
	return repository.storeItem(getKey(note.Id), note)
}

func (repository *FilesystemItemRepository) GetItems() []entities.Manageable {
	items := loadItems(repository.StorageDir)
	result := make([]entities.Manageable, len(items))

	for _, item := range items {
		index := entities.GetId(item) - 1
		if item != nil {
			result[index] = entities.CreateItem(item)
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

func (repository *FilesystemItemRepository) UpdateItem(id int, title *string, description *string, tags ...string) error {
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
	archivedItems := loadItems(repository.StorageDir)
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
